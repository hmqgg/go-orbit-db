package orbitdb

import (
	"context"

	coreapi "github.com/ipfs/interface-go-ipfs-core"
	"github.com/pkg/errors"

	"github.com/hmqgg/go-orbit-db/accesscontroller/ipfs"
	"github.com/hmqgg/go-orbit-db/accesscontroller/orbitdb"
	"github.com/hmqgg/go-orbit-db/accesscontroller/simple"
	"github.com/hmqgg/go-orbit-db/baseorbitdb"
	"github.com/hmqgg/go-orbit-db/iface"
	"github.com/hmqgg/go-orbit-db/stores/documentstore"
	"github.com/hmqgg/go-orbit-db/stores/eventlogstore"
	"github.com/hmqgg/go-orbit-db/stores/kvstore"
)

type orbitDB struct {
	baseorbitdb.BaseOrbitDB
}

// OrbitDB An alias of the type defined in the iface package
type OrbitDB = iface.OrbitDB

// Store An alias of the type defined in the iface package
type Store = iface.Store

// EventLogStore An alias of the type defined in the iface package
type EventLogStore = iface.EventLogStore

// KeyValueStore An alias of the type defined in the iface package
type KeyValueStore = iface.KeyValueStore

// DocumentStore An alias of the type defined in the iface package
type DocumentStore = iface.DocumentStore

// StoreIndex An alias of the type defined in the iface package
type StoreIndex = iface.StoreIndex

// StoreConstructor An alias of the type defined in the iface package
type StoreConstructor = iface.StoreConstructor

// IndexConstructor An alias of the type defined in the iface package
type IndexConstructor = iface.IndexConstructor

// OnWritePrototype An alias of the type defined in the iface package
type OnWritePrototype = iface.OnWritePrototype

// StreamOptions An alias of the type defined in the iface package
type StreamOptions = iface.StreamOptions

// CreateDBOptions An alias of the type defined in the iface package
type CreateDBOptions = iface.CreateDBOptions

// CreateDocumentDBOptions An alias of the type defined in the iface package
type CreateDocumentDBOptions = iface.CreateDocumentDBOptions

// DetermineAddressOptions An alias of the type defined in the iface package
type DetermineAddressOptions = iface.DetermineAddressOptions

// NewOrbitDBOptions Options for a new OrbitDB instance
type NewOrbitDBOptions = baseorbitdb.NewOrbitDBOptions

// NewOrbitDB Creates a new OrbitDB instance with default access controllers and store types
func NewOrbitDB(ctx context.Context, i coreapi.CoreAPI, options *NewOrbitDBOptions) (iface.OrbitDB, error) {
	odb, err := baseorbitdb.NewOrbitDB(ctx, i, options)

	if err != nil {
		return nil, err
	}

	odb.RegisterStoreType("eventlog", eventlogstore.NewOrbitDBEventLogStore)
	odb.RegisterStoreType("keyvalue", kvstore.NewOrbitDBKeyValue)
	odb.RegisterStoreType("docstore", documentstore.NewOrbitDBDocumentStore)

	_ = odb.RegisterAccessControllerType(ipfs.NewIPFSAccessController)
	_ = odb.RegisterAccessControllerType(orbitdb.NewOrbitDBAccessController)
	_ = odb.RegisterAccessControllerType(simple.NewSimpleAccessController)

	return &orbitDB{
		BaseOrbitDB: odb,
	}, nil
}

func (o *orbitDB) Log(ctx context.Context, address string, options *CreateDBOptions) (EventLogStore, error) {
	if options == nil {
		options = &CreateDBOptions{}
	}

	options.Create = boolPtr(true)
	options.StoreType = stringPtr("eventlog")
	store, err := o.Open(ctx, address, options)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open database")
	}

	logStore, ok := store.(EventLogStore)
	if !ok {
		return nil, errors.New("unable to cast store to log")
	}

	return logStore, nil
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func (o *orbitDB) KeyValue(ctx context.Context, address string, options *CreateDBOptions) (KeyValueStore, error) {
	if options == nil {
		options = &CreateDBOptions{}
	}

	options.Create = boolPtr(true)
	options.StoreType = stringPtr("keyvalue")

	store, err := o.Open(ctx, address, options)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open database")
	}

	kvStore, ok := store.(KeyValueStore)
	if !ok {
		return nil, errors.New("unable to cast store to keyvalue")
	}

	return kvStore, nil
}

func (o *orbitDB) Docs(ctx context.Context, address string, options *CreateDBOptions) (DocumentStore, error) {
	if options == nil {
		options = &CreateDBOptions{}
	}

	options.Create = boolPtr(true)
	options.StoreType = stringPtr("docstore")

	store, err := o.Open(ctx, address, options)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open database")
	}

	documentStore, ok := store.(DocumentStore)
	if !ok {
		return nil, errors.New("unable to cast store to document")
	}

	return documentStore, nil
}

var _ OrbitDB = (*orbitDB)(nil)
