package dht

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	cid "gx/ipfs/QmPSQnBKM9g7BaUcZCvswUJVscQ1ipjmwxN5PXCjkp9EQ7/go-cid"
	u "gx/ipfs/QmPdKqUcHGFdeSpvjVoaTRPPstGif9GBZb5Q56RVw9o69A/go-ipfs-util"
	ob "gx/ipfs/QmPpYHPRGVpSJTkQDQDwTYZ1cYUR2NM4HS6M3iAXi8aoUa/go-libp2p-kad-dht/ob"
	pb "gx/ipfs/QmPpYHPRGVpSJTkQDQDwTYZ1cYUR2NM4HS6M3iAXi8aoUa/go-libp2p-kad-dht/pb"
	peer "gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"
	pstore "gx/ipfs/QmTTJcDL3gsnGDALjh2fDGg1onGRUdVgNL2hU2WEZcVrMX/go-libp2p-peerstore"
	inet "gx/ipfs/QmXuRkCR7BNQa9uqfpTiFWsTQLzmTWYg91Ja1w95gnqb6u/go-libp2p-net"
	recpb "gx/ipfs/Qma9Eqp16mNHDX1EL73pcxhFfzbyXVcAYtaDd1xdmDRDtL/go-libp2p-record/pb"
	ds "gx/ipfs/QmaRb5yNXKonhbkpNxNawoydk4N6es6b4fPj19sjEKsh5D/go-datastore"
	proto "gx/ipfs/QmdxUuburamoF6zF9qjeQC4WYcWGbWuRmdLacMEsW8ioD8/gogo-protobuf/proto"
	base32 "gx/ipfs/QmfVj3x4D6Jkq9SEoi5n2NmoUomLwoeiwnYz2KQa15wRw6/base32"
)

// The number of closer peers to send on requests.
var CloserPeerCount = KValue

// dhthandler specifies the signature of functions that handle DHT messages.
type dhtHandler func(context.Context, peer.ID, *pb.Message) (*pb.Message, error)

func (dht *IpfsDHT) handlerForMsgType(t pb.Message_MessageType) dhtHandler {
	switch t {
	case pb.Message_GET_VALUE:
		return dht.handleGetValue
	case pb.Message_PUT_VALUE:
		return dht.handlePutValue
	case pb.Message_FIND_NODE:
		return dht.handleFindPeer
	case pb.Message_ADD_PROVIDER:
		return dht.handleAddProvider
	case pb.Message_GET_PROVIDERS:
		return dht.handleGetProviders
	case pb.Message_PING:
		return dht.handlePing
	default:
		return nil
	}
}

func (dht *IpfsDHT) handleGetValue(ctx context.Context, p peer.ID, pmes *pb.Message) (_ *pb.Message, err error) {
	ctx = log.Start(ctx, "handleGetValue")
	log.SetTag(ctx, "peer", p)
	defer func() { log.FinishWithErr(ctx, err) }()
	log.Debugf("%s handleGetValue for key: %s", dht.self, pmes.GetKey())

	// setup response
	resp := pb.NewMessage(pmes.GetType(), pmes.GetKey(), pmes.GetClusterLevel())

	// first, is there even a key?
	k := pmes.GetKey()
	if len(k) == 0 {
		return nil, errors.New("handleGetValue but no key was provided")
		// TODO: send back an error response? could be bad, but the other node's hanging.
	}

	rec, err := dht.checkLocalDatastore(k)
	if err != nil {
		return nil, err
	}
	resp.Record = rec

	// Find closest peer on given cluster to desired key and reply with that info
	closer := dht.betterPeersToQuery(pmes, p, CloserPeerCount)
	if len(closer) > 0 {
		closerinfos := pstore.PeerInfos(dht.peerstore, closer)
		for _, pi := range closerinfos {
			log.Debugf("handleGetValue returning closer peer: '%s'", pi.ID)
			if len(pi.Addrs) < 1 {
				log.Warningf(`no addresses on peer being sent!
					[local:%s]
					[sending:%s]
					[remote:%s]`, dht.self, pi.ID, p)
			}
		}

		resp.CloserPeers = pb.PeerInfosToPBPeers(dht.host.Network(), closerinfos)
	}

	return resp, nil
}

func (dht *IpfsDHT) checkLocalDatastore(k []byte) (*recpb.Record, error) {
	log.Debugf("%s handleGetValue looking into ds", dht.self)
	dskey := convertToDsKey(k)
	buf, err := dht.datastore.Get(dskey)
	log.Debugf("%s handleGetValue looking into ds GOT %v", dht.self, buf)

	if err == ds.ErrNotFound {
		return nil, nil
	}

	// if we got an unexpected error, bail.
	if err != nil {
		return nil, err
	}

	// if we have the value, send it back
	log.Debugf("%s handleGetValue success!", dht.self)

	rec := new(recpb.Record)
	err = proto.Unmarshal(buf, rec)
	if err != nil {
		log.Debug("failed to unmarshal DHT record from datastore")
		return nil, err
	}

	var recordIsBad bool
	recvtime, err := u.ParseRFC3339(rec.GetTimeReceived())
	if err != nil {
		log.Info("either no receive time set on record, or it was invalid: ", err)
		recordIsBad = true
	}

	if time.Now().Sub(recvtime) > MaxRecordAge {
		log.Debug("old record found, tossing.")
		recordIsBad = true
	}

	// NOTE: We do not verify the record here beyond checking these timestamps.
	// we put the burden of checking the records on the requester as checking a record
	// may be computationally expensive

	if recordIsBad {
		err := dht.datastore.Delete(dskey)
		if err != nil {
			log.Error("Failed to delete bad record from datastore: ", err)
		}

		return nil, nil // can treat this as not having the record at all
	}

	return rec, nil
}

// Cleans the record (to avoid storing arbitrary data).
func cleanRecord(rec *recpb.Record) {
	rec.TimeReceived = ""
}

// Store a value in this peer local storage
func (dht *IpfsDHT) handlePutValue(ctx context.Context, p peer.ID, pmes *pb.Message) (_ *pb.Message, err error) {
	ctx = log.Start(ctx, "handlePutValue")
	log.SetTag(ctx, "peer", p)
	defer func() { log.FinishWithErr(ctx, err) }()

	rec := pmes.GetRecord()
	if rec == nil {
		log.Infof("Got nil record from: %s", p.Pretty())
		return nil, errors.New("nil record")
	}

	if !bytes.Equal(pmes.GetKey(), rec.GetKey()) {
		return nil, errors.New("put key doesn't match record key")
	}

	cleanRecord(rec)

	// Make sure the record is valid (not expired, valid signature etc)
	if err = dht.Validator.Validate(string(rec.GetKey()), rec.GetValue()); err != nil {
		log.Warningf("Bad dht record in PUT from: %s. %s", p.Pretty(), err)
		return nil, err
	}

	dskey := convertToDsKey(rec.GetKey())

	// Make sure the new record is "better" than the record we have locally.
	// This prevents a record with for example a lower sequence number from
	// overwriting a record with a higher sequence number.
	existing, err := dht.getRecordFromDatastore(dskey)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		recs := [][]byte{rec.GetValue(), existing.GetValue()}
		i, err := dht.Validator.Select(string(rec.GetKey()), recs)
		if err != nil {
			log.Warningf("Bad dht record in PUT from %s: %s", p.Pretty(), err)
			return nil, err
		}
		if i != 0 {
			log.Infof("DHT record in PUT from %s is older than existing record. Ignoring", p.Pretty())
			return nil, errors.New("old record")
		}
	}

	// record the time we receive every record
	rec.TimeReceived = u.FormatRFC3339(time.Now())

	data, err := proto.Marshal(rec)
	if err != nil {
		return nil, err
	}

	err = dht.datastore.Put(dskey, data)
	log.Debugf("%s handlePutValue %v", dht.self, dskey)
	return pmes, err
}

// returns nil, nil when either nothing is found or the value found doesn't properly validate.
// returns nil, some_error when there's a *datastore* error (i.e., something goes very wrong)
func (dht *IpfsDHT) getRecordFromDatastore(dskey ds.Key) (*recpb.Record, error) {
	buf, err := dht.datastore.Get(dskey)
	if err == ds.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		log.Errorf("Got error retrieving record with key %s from datastore: %s", dskey, err)
		return nil, err
	}
	rec := new(recpb.Record)
	err = proto.Unmarshal(buf, rec)
	if err != nil {
		// Bad data in datastore, log it but don't return an error, we'll just overwrite it
		log.Errorf("Bad record data stored in datastore with key %s: could not unmarshal record", dskey)
		return nil, nil
	}

	err = dht.Validator.Validate(string(rec.GetKey()), rec.GetValue())
	if err != nil {
		// Invalid record in datastore, probably expired but don't return an error,
		// we'll just overwrite it
		log.Debugf("Local record verify failed: %s (discarded)", err)
		return nil, nil
	}

	return rec, nil
}

func (dht *IpfsDHT) handlePing(_ context.Context, p peer.ID, pmes *pb.Message) (*pb.Message, error) {
	log.Debugf("%s Responding to ping from %s!\n", dht.self, p)
	return pmes, nil
}

func (dht *IpfsDHT) handleFindPeer(ctx context.Context, p peer.ID, pmes *pb.Message) (_ *pb.Message, _err error) {
	ctx = log.Start(ctx, "handleFindPeer")
	defer func() { log.FinishWithErr(ctx, _err) }()
	log.SetTag(ctx, "peer", p)
	resp := pb.NewMessage(pmes.GetType(), nil, pmes.GetClusterLevel())
	var closest []peer.ID

	// if looking for self... special case where we send it on CloserPeers.
	targetPid := peer.ID(pmes.GetKey())
	if targetPid == dht.self {
		closest = []peer.ID{dht.self}
	} else {
		closest = dht.betterPeersToQuery(pmes, p, CloserPeerCount)

		// Never tell a peer about itself.
		if targetPid != p {
			// If we're connected to the target peer, report their
			// peer info. This makes FindPeer work even if the
			// target peer isn't in our routing table.
			//
			// Alternatively, we could just check our peerstore.
			// However, we don't want to return out of date
			// information. We can change this in the future when we
			// add a progressive, asynchronous `SearchPeer` function
			// and improve peer routing in the host.
			switch dht.host.Network().Connectedness(targetPid) {
			case inet.Connected, inet.CanConnect:
				closest = append(closest, targetPid)
			}
		}
	}

	if closest == nil {
		log.Infof("%s handleFindPeer %s: could not find anything.", dht.self, p)
		return resp, nil
	}

	closestinfos := pstore.PeerInfos(dht.peerstore, closest)
	// possibly an over-allocation but this array is temporary anyways.
	withAddresses := make([]pstore.PeerInfo, 0, len(closestinfos))
	for _, pi := range closestinfos {
		if len(pi.Addrs) > 0 {
			withAddresses = append(withAddresses, pi)
			log.Debugf("handleFindPeer: sending back '%s'", pi.ID)
		}
	}

	resp.CloserPeers = pb.PeerInfosToPBPeers(dht.host.Network(), withAddresses)
	return resp, nil
}

func (dht *IpfsDHT) handleGetProviders(ctx context.Context, p peer.ID, pmes *pb.Message) (_ *pb.Message, _err error) {
	ctx = log.Start(ctx, "handleGetProviders")
	defer func() { log.FinishWithErr(ctx, _err) }()
	log.SetTag(ctx, "peer", p)

	resp := pb.NewMessage(pmes.GetType(), pmes.GetKey(), pmes.GetClusterLevel())
	c, err := cid.Cast([]byte(pmes.GetKey()))
	if err != nil {
		return nil, err
	}
	log.SetTag(ctx, "key", c)

	// debug logging niceness.
	reqDesc := fmt.Sprintf("%s handleGetProviders(%s, %s): ", dht.self, p, c)
	log.Debugf("%s begin", reqDesc)
	defer log.Debugf("%s end", reqDesc)

	// check if we have this value, to add ourselves as provider.
	has, err := dht.datastore.Has(convertToDsKey(c.Bytes()))
	if err != nil && err != ds.ErrNotFound {
		log.Debugf("unexpected datastore error: %v\n", err)
		has = false
	}

	// setup providers
	providers := dht.providers.GetProviders(ctx, c)
	if has {
		providers = append(providers, dht.self)
		log.Debugf("%s have the value. added self as provider", reqDesc)
	}

	if providers != nil && len(providers) > 0 {
		infos := pstore.PeerInfos(dht.peerstore, providers)
		resp.ProviderPeers = pb.PeerInfosToPBPeers(dht.host.Network(), infos)
		log.Debugf("%s have %d providers: %s", reqDesc, len(providers), infos)
	}

	// Also send closer peers.
	closer := dht.betterPeersToQuery(pmes, p, CloserPeerCount)
	if closer != nil {
		infos := pstore.PeerInfos(dht.peerstore, closer)
		resp.CloserPeers = pb.PeerInfosToPBPeers(dht.host.Network(), infos)
		log.Debugf("%s have %d closer peers: %s", reqDesc, len(closer), infos)
	}

	return resp, nil
}

func (dht *IpfsDHT) handleAddProvider(ctx context.Context, p peer.ID, pmes *pb.Message) (_ *pb.Message, _err error) {
	ctx = log.Start(ctx, "handleAddProvider")
	defer func() { log.FinishWithErr(ctx, _err) }()
	log.SetTag(ctx, "peer", p)

	c, err := cid.Cast([]byte(pmes.GetKey()))
	if err != nil {
		return nil, err
	}
	log.SetTag(ctx, "key", c)

	log.Debugf("%s adding %s as a provider for '%s'\n", dht.self, p, c)

	// add provider should use the address given in the message
	pinfos := pb.PBPeersToPeerInfos(pmes.GetProviderPeers())
	for _, pi := range pinfos {
		// OpenBazaar: there was a line here that checks to make sure the peerID equals the provider ID and
		// returns an error if it does not. This line has been removed since we want to be able to store
		// our `pointers` in the DHT as a provider which does not have an ID matching the peerID.

		if len(pi.Addrs) < 1 {
			log.Debugf("%s got no valid addresses for provider %s. Ignore.", dht.self, p)
			continue
		}

		// OpenBazaar: if this is a pointer we want to use the PointerAddrTTL instead of the normal
		// ProviderAddrTTL.
		ttl := pstore.ProviderAddrTTL
		if ob.IsPointer(pi.ID) {
			ttl = ob.PointerAddrTTL
		}

		log.Infof("received provider %s for %s (addrs: %s)", p, c, pi.Addrs)
		if pi.ID != dht.self { // don't add own addrs.
			// add the received addresses to our peerstore.
			// OpenBazaar: make sure to use `ttl` as a parameter here.
			dht.peerstore.AddAddrs(pi.ID, pi.Addrs, ttl)
		}
		// OpenBazaar: the line below originally used p instead of pi.ID as the peerID. For normal providers
		// this would be the same value and it doesn't matter which one is used. But for our `pointers` p is
		// not the value we want to add, but rather pi.ID. So we substitute pi.ID here.
		dht.providers.AddProvider(ctx, c, pi.ID)
	}

	return nil, nil
}

func convertToDsKey(s []byte) ds.Key {
	return ds.NewKey(base32.RawStdEncoding.EncodeToString(s))
}
