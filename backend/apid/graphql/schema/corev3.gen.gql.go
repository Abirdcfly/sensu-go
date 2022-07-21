// Code generated by scripts/gengraphql.go. DO NOT EDIT.

package schema

import (
	errors "errors"
	graphql1 "github.com/graphql-go/graphql"
	graphql "github.com/sensu/sensu-go/graphql"
)

//
// CoreV3EntityConfigFieldResolvers represents a collection of methods whose products represent the
// response values of the 'CoreV3EntityConfig' type.
type CoreV3EntityConfigFieldResolvers interface {
	// Metadata implements response to request for 'metadata' field.
	Metadata(p graphql.ResolveParams) (interface{}, error)

	// Entity_class implements response to request for 'entity_class' field.
	Entity_class(p graphql.ResolveParams) (string, error)

	// User implements response to request for 'user' field.
	User(p graphql.ResolveParams) (string, error)

	// Subscriptions implements response to request for 'subscriptions' field.
	Subscriptions(p graphql.ResolveParams) ([]string, error)

	// Deregister implements response to request for 'deregister' field.
	Deregister(p graphql.ResolveParams) (bool, error)

	// Deregistration implements response to request for 'deregistration' field.
	Deregistration(p graphql.ResolveParams) (interface{}, error)

	// Keepalive_handlers implements response to request for 'keepalive_handlers' field.
	Keepalive_handlers(p graphql.ResolveParams) ([]string, error)

	// Redact implements response to request for 'redact' field.
	Redact(p graphql.ResolveParams) ([]string, error)
}

// CoreV3EntityConfigAliases implements all methods on CoreV3EntityConfigFieldResolvers interface by using reflection to
// match name of field to a field on the given value. Intent is reduce friction
// of writing new resolvers by removing all the instances where you would simply
// have the resolvers method return a field.
type CoreV3EntityConfigAliases struct{}

// Metadata implements response to request for 'metadata' field.
func (_ CoreV3EntityConfigAliases) Metadata(p graphql.ResolveParams) (interface{}, error) {
	val, err := graphql.DefaultResolver(p.Source, p.Info.FieldName)
	return val, err
}

// Entity_class implements response to request for 'entity_class' field.
func (_ CoreV3EntityConfigAliases) Entity_class(p graphql.ResolveParams) (string, error) {
	val, err := graphql.DefaultResolver(p.Source, p.Info.FieldName)
	ret, ok := val.(string)
	if err != nil {
		return ret, err
	}
	if !ok {
		return ret, errors.New("unable to coerce value for field 'entity_class'")
	}
	return ret, err
}

// User implements response to request for 'user' field.
func (_ CoreV3EntityConfigAliases) User(p graphql.ResolveParams) (string, error) {
	val, err := graphql.DefaultResolver(p.Source, p.Info.FieldName)
	ret, ok := val.(string)
	if err != nil {
		return ret, err
	}
	if !ok {
		return ret, errors.New("unable to coerce value for field 'user'")
	}
	return ret, err
}

// Subscriptions implements response to request for 'subscriptions' field.
func (_ CoreV3EntityConfigAliases) Subscriptions(p graphql.ResolveParams) ([]string, error) {
	val, err := graphql.DefaultResolver(p.Source, p.Info.FieldName)
	ret, ok := val.([]string)
	if err != nil {
		return ret, err
	}
	if !ok {
		return ret, errors.New("unable to coerce value for field 'subscriptions'")
	}
	return ret, err
}

// Deregister implements response to request for 'deregister' field.
func (_ CoreV3EntityConfigAliases) Deregister(p graphql.ResolveParams) (bool, error) {
	val, err := graphql.DefaultResolver(p.Source, p.Info.FieldName)
	ret, ok := val.(bool)
	if err != nil {
		return ret, err
	}
	if !ok {
		return ret, errors.New("unable to coerce value for field 'deregister'")
	}
	return ret, err
}

// Deregistration implements response to request for 'deregistration' field.
func (_ CoreV3EntityConfigAliases) Deregistration(p graphql.ResolveParams) (interface{}, error) {
	val, err := graphql.DefaultResolver(p.Source, p.Info.FieldName)
	return val, err
}

// Keepalive_handlers implements response to request for 'keepalive_handlers' field.
func (_ CoreV3EntityConfigAliases) Keepalive_handlers(p graphql.ResolveParams) ([]string, error) {
	val, err := graphql.DefaultResolver(p.Source, p.Info.FieldName)
	ret, ok := val.([]string)
	if err != nil {
		return ret, err
	}
	if !ok {
		return ret, errors.New("unable to coerce value for field 'keepalive_handlers'")
	}
	return ret, err
}

// Redact implements response to request for 'redact' field.
func (_ CoreV3EntityConfigAliases) Redact(p graphql.ResolveParams) ([]string, error) {
	val, err := graphql.DefaultResolver(p.Source, p.Info.FieldName)
	ret, ok := val.([]string)
	if err != nil {
		return ret, err
	}
	if !ok {
		return ret, errors.New("unable to coerce value for field 'redact'")
	}
	return ret, err
}

// CoreV3EntityConfigType EntityConfig represents entity configuration.
var CoreV3EntityConfigType = graphql.NewType("CoreV3EntityConfig", graphql.ObjectKind)

// RegisterCoreV3EntityConfig registers CoreV3EntityConfig object type with given service.
func RegisterCoreV3EntityConfig(svc *graphql.Service, impl CoreV3EntityConfigFieldResolvers) {
	svc.RegisterObject(_ObjectTypeCoreV3EntityConfigDesc, impl)
}
func _ObjTypeCoreV3EntityConfigMetadataHandler(impl interface{}) graphql1.FieldResolveFn {
	resolver := impl.(interface {
		Metadata(p graphql.ResolveParams) (interface{}, error)
	})
	return func(frp graphql1.ResolveParams) (interface{}, error) {
		return resolver.Metadata(frp)
	}
}

func _ObjTypeCoreV3EntityConfigEntity_classHandler(impl interface{}) graphql1.FieldResolveFn {
	resolver := impl.(interface {
		Entity_class(p graphql.ResolveParams) (string, error)
	})
	return func(frp graphql1.ResolveParams) (interface{}, error) {
		return resolver.Entity_class(frp)
	}
}

func _ObjTypeCoreV3EntityConfigUserHandler(impl interface{}) graphql1.FieldResolveFn {
	resolver := impl.(interface {
		User(p graphql.ResolveParams) (string, error)
	})
	return func(frp graphql1.ResolveParams) (interface{}, error) {
		return resolver.User(frp)
	}
}

func _ObjTypeCoreV3EntityConfigSubscriptionsHandler(impl interface{}) graphql1.FieldResolveFn {
	resolver := impl.(interface {
		Subscriptions(p graphql.ResolveParams) ([]string, error)
	})
	return func(frp graphql1.ResolveParams) (interface{}, error) {
		return resolver.Subscriptions(frp)
	}
}

func _ObjTypeCoreV3EntityConfigDeregisterHandler(impl interface{}) graphql1.FieldResolveFn {
	resolver := impl.(interface {
		Deregister(p graphql.ResolveParams) (bool, error)
	})
	return func(frp graphql1.ResolveParams) (interface{}, error) {
		return resolver.Deregister(frp)
	}
}

func _ObjTypeCoreV3EntityConfigDeregistrationHandler(impl interface{}) graphql1.FieldResolveFn {
	resolver := impl.(interface {
		Deregistration(p graphql.ResolveParams) (interface{}, error)
	})
	return func(frp graphql1.ResolveParams) (interface{}, error) {
		return resolver.Deregistration(frp)
	}
}

func _ObjTypeCoreV3EntityConfigKeepalive_handlersHandler(impl interface{}) graphql1.FieldResolveFn {
	resolver := impl.(interface {
		Keepalive_handlers(p graphql.ResolveParams) ([]string, error)
	})
	return func(frp graphql1.ResolveParams) (interface{}, error) {
		return resolver.Keepalive_handlers(frp)
	}
}

func _ObjTypeCoreV3EntityConfigRedactHandler(impl interface{}) graphql1.FieldResolveFn {
	resolver := impl.(interface {
		Redact(p graphql.ResolveParams) ([]string, error)
	})
	return func(frp graphql1.ResolveParams) (interface{}, error) {
		return resolver.Redact(frp)
	}
}

func _ObjectTypeCoreV3EntityConfigConfigFn() graphql1.ObjectConfig {
	return graphql1.ObjectConfig{
		Description: "EntityConfig represents entity configuration.",
		Fields: graphql1.Fields{
			"deregister": &graphql1.Field{
				Args:              graphql1.FieldConfigArgument{},
				DeprecationReason: "",
				Description:       "Deregister, if true, will result in the entity being deleted when the\nentity is an agent, and the agent disconnects its session.",
				Name:              "deregister",
				Type:              graphql1.NewNonNull(graphql1.Boolean),
			},
			"deregistration": &graphql1.Field{
				Args:              graphql1.FieldConfigArgument{},
				DeprecationReason: "",
				Description:       "Deregistration contains configuration for Sensu entity de-registration.",
				Name:              "deregistration",
				Type:              graphql1.NewNonNull(graphql.OutputType("CoreV2Deregistration")),
			},
			"entity_class": &graphql1.Field{
				Args:              graphql1.FieldConfigArgument{},
				DeprecationReason: "",
				Description:       "EntityClass represents the class of the entity. It can be \"agent\",\n\"proxy\", or \"backend\".",
				Name:              "entity_class",
				Type:              graphql1.NewNonNull(graphql1.String),
			},
			"keepalive_handlers": &graphql1.Field{
				Args:              graphql1.FieldConfigArgument{},
				DeprecationReason: "",
				Description:       "KeepaliveHandlers contains a list of handlers to use for the entity's\nkeepalive events.",
				Name:              "keepalive_handlers",
				Type:              graphql1.NewNonNull(graphql1.NewList(graphql1.NewNonNull(graphql1.String))),
			},
			"metadata": &graphql1.Field{
				Args:              graphql1.FieldConfigArgument{},
				DeprecationReason: "",
				Description:       "Metadata contains the name, namespace, labels and annotations of the\nentity.",
				Name:              "metadata",
				Type:              graphql.OutputType("ObjectMeta"),
			},
			"redact": &graphql1.Field{
				Args:              graphql1.FieldConfigArgument{},
				DeprecationReason: "",
				Description:       "Redact contains the fields to redact on the entity, if the entity is an]\nagent entity.",
				Name:              "redact",
				Type:              graphql1.NewNonNull(graphql1.NewList(graphql1.NewNonNull(graphql1.String))),
			},
			"subscriptions": &graphql1.Field{
				Args:              graphql1.FieldConfigArgument{},
				DeprecationReason: "",
				Description:       "Subscriptions are a weak relationship between entities and checks. The\nscheduler uses subscriptions to make entities to checks when scheduling.",
				Name:              "subscriptions",
				Type:              graphql1.NewNonNull(graphql1.NewList(graphql1.NewNonNull(graphql1.String))),
			},
			"user": &graphql1.Field{
				Args:              graphql1.FieldConfigArgument{},
				DeprecationReason: "",
				Description:       "User is the username the entity is connecting as, if the entity is an\nagent entity.",
				Name:              "user",
				Type:              graphql1.NewNonNull(graphql1.String),
			},
		},
		Interfaces: []*graphql1.Interface{},
		IsTypeOf: func(_ graphql1.IsTypeOfParams) bool {
			// NOTE:
			// Panic by default. Intent is that when Service is invoked, values of
			// these fields are updated with instantiated resolvers. If these
			// defaults are called it is most certainly programmer err.
			// If you're see this comment then: 'Whoops! Sorry, my bad.'
			panic("Unimplemented; see CoreV3EntityConfigFieldResolvers.")
		},
		Name: "CoreV3EntityConfig",
	}
}

// describe CoreV3EntityConfig's configuration; kept private to avoid unintentional tampering of configuration at runtime.
var _ObjectTypeCoreV3EntityConfigDesc = graphql.ObjectDesc{
	Config: _ObjectTypeCoreV3EntityConfigConfigFn,
	FieldHandlers: map[string]graphql.FieldHandler{
		"deregister":         _ObjTypeCoreV3EntityConfigDeregisterHandler,
		"deregistration":     _ObjTypeCoreV3EntityConfigDeregistrationHandler,
		"entity_class":       _ObjTypeCoreV3EntityConfigEntity_classHandler,
		"keepalive_handlers": _ObjTypeCoreV3EntityConfigKeepalive_handlersHandler,
		"metadata":           _ObjTypeCoreV3EntityConfigMetadataHandler,
		"redact":             _ObjTypeCoreV3EntityConfigRedactHandler,
		"subscriptions":      _ObjTypeCoreV3EntityConfigSubscriptionsHandler,
		"user":               _ObjTypeCoreV3EntityConfigUserHandler,
	},
}

//
// CoreV3EntityStateFieldResolvers represents a collection of methods whose products represent the
// response values of the 'CoreV3EntityState' type.
type CoreV3EntityStateFieldResolvers interface {
	// Metadata implements response to request for 'metadata' field.
	Metadata(p graphql.ResolveParams) (interface{}, error)

	// System implements response to request for 'system' field.
	System(p graphql.ResolveParams) (interface{}, error)

	// Last_seen implements response to request for 'last_seen' field.
	Last_seen(p graphql.ResolveParams) (int, error)

	// Sensu_agent_version implements response to request for 'sensu_agent_version' field.
	Sensu_agent_version(p graphql.ResolveParams) (string, error)
}

// CoreV3EntityStateAliases implements all methods on CoreV3EntityStateFieldResolvers interface by using reflection to
// match name of field to a field on the given value. Intent is reduce friction
// of writing new resolvers by removing all the instances where you would simply
// have the resolvers method return a field.
type CoreV3EntityStateAliases struct{}

// Metadata implements response to request for 'metadata' field.
func (_ CoreV3EntityStateAliases) Metadata(p graphql.ResolveParams) (interface{}, error) {
	val, err := graphql.DefaultResolver(p.Source, p.Info.FieldName)
	return val, err
}

// System implements response to request for 'system' field.
func (_ CoreV3EntityStateAliases) System(p graphql.ResolveParams) (interface{}, error) {
	val, err := graphql.DefaultResolver(p.Source, p.Info.FieldName)
	return val, err
}

// Last_seen implements response to request for 'last_seen' field.
func (_ CoreV3EntityStateAliases) Last_seen(p graphql.ResolveParams) (int, error) {
	val, err := graphql.DefaultResolver(p.Source, p.Info.FieldName)
	ret, ok := graphql1.Int.ParseValue(val).(int)
	if err != nil {
		return ret, err
	}
	if !ok {
		return ret, errors.New("unable to coerce value for field 'last_seen'")
	}
	return ret, err
}

// Sensu_agent_version implements response to request for 'sensu_agent_version' field.
func (_ CoreV3EntityStateAliases) Sensu_agent_version(p graphql.ResolveParams) (string, error) {
	val, err := graphql.DefaultResolver(p.Source, p.Info.FieldName)
	ret, ok := val.(string)
	if err != nil {
		return ret, err
	}
	if !ok {
		return ret, errors.New("unable to coerce value for field 'sensu_agent_version'")
	}
	return ret, err
}

/*
CoreV3EntityStateType EntityState represents entity state. Unlike configuration, state is
typically only maintained for agent entities, although it can be maintained
for proxy entities in certain circumstances.
*/
var CoreV3EntityStateType = graphql.NewType("CoreV3EntityState", graphql.ObjectKind)

// RegisterCoreV3EntityState registers CoreV3EntityState object type with given service.
func RegisterCoreV3EntityState(svc *graphql.Service, impl CoreV3EntityStateFieldResolvers) {
	svc.RegisterObject(_ObjectTypeCoreV3EntityStateDesc, impl)
}
func _ObjTypeCoreV3EntityStateMetadataHandler(impl interface{}) graphql1.FieldResolveFn {
	resolver := impl.(interface {
		Metadata(p graphql.ResolveParams) (interface{}, error)
	})
	return func(frp graphql1.ResolveParams) (interface{}, error) {
		return resolver.Metadata(frp)
	}
}

func _ObjTypeCoreV3EntityStateSystemHandler(impl interface{}) graphql1.FieldResolveFn {
	resolver := impl.(interface {
		System(p graphql.ResolveParams) (interface{}, error)
	})
	return func(frp graphql1.ResolveParams) (interface{}, error) {
		return resolver.System(frp)
	}
}

func _ObjTypeCoreV3EntityStateLast_seenHandler(impl interface{}) graphql1.FieldResolveFn {
	resolver := impl.(interface {
		Last_seen(p graphql.ResolveParams) (int, error)
	})
	return func(frp graphql1.ResolveParams) (interface{}, error) {
		return resolver.Last_seen(frp)
	}
}

func _ObjTypeCoreV3EntityStateSensu_agent_versionHandler(impl interface{}) graphql1.FieldResolveFn {
	resolver := impl.(interface {
		Sensu_agent_version(p graphql.ResolveParams) (string, error)
	})
	return func(frp graphql1.ResolveParams) (interface{}, error) {
		return resolver.Sensu_agent_version(frp)
	}
}

func _ObjectTypeCoreV3EntityStateConfigFn() graphql1.ObjectConfig {
	return graphql1.ObjectConfig{
		Description: "EntityState represents entity state. Unlike configuration, state is\ntypically only maintained for agent entities, although it can be maintained\nfor proxy entities in certain circumstances.",
		Fields: graphql1.Fields{
			"last_seen": &graphql1.Field{
				Args:              graphql1.FieldConfigArgument{},
				DeprecationReason: "",
				Description:       "LastSeen is a unix timestamp that represents when the entity was last\nobserved by the keepalive system.",
				Name:              "last_seen",
				Type:              graphql1.NewNonNull(graphql1.Int),
			},
			"metadata": &graphql1.Field{
				Args:              graphql1.FieldConfigArgument{},
				DeprecationReason: "",
				Description:       "Metadata contains the name, namespace, labels and annotations of the\nentity.",
				Name:              "metadata",
				Type:              graphql.OutputType("ObjectMeta"),
			},
			"sensu_agent_version": &graphql1.Field{
				Args:              graphql1.FieldConfigArgument{},
				DeprecationReason: "",
				Description:       "SensuAgentVersion is the sensu-agent version.",
				Name:              "sensu_agent_version",
				Type:              graphql1.NewNonNull(graphql1.String),
			},
			"system": &graphql1.Field{
				Args:              graphql1.FieldConfigArgument{},
				DeprecationReason: "",
				Description:       "System contains information about the system that the Agent process\nis running on, used for additional Entity context.",
				Name:              "system",
				Type:              graphql1.NewNonNull(graphql.OutputType("CoreV2System")),
			},
		},
		Interfaces: []*graphql1.Interface{},
		IsTypeOf: func(_ graphql1.IsTypeOfParams) bool {
			// NOTE:
			// Panic by default. Intent is that when Service is invoked, values of
			// these fields are updated with instantiated resolvers. If these
			// defaults are called it is most certainly programmer err.
			// If you're see this comment then: 'Whoops! Sorry, my bad.'
			panic("Unimplemented; see CoreV3EntityStateFieldResolvers.")
		},
		Name: "CoreV3EntityState",
	}
}

// describe CoreV3EntityState's configuration; kept private to avoid unintentional tampering of configuration at runtime.
var _ObjectTypeCoreV3EntityStateDesc = graphql.ObjectDesc{
	Config: _ObjectTypeCoreV3EntityStateConfigFn,
	FieldHandlers: map[string]graphql.FieldHandler{
		"last_seen":           _ObjTypeCoreV3EntityStateLast_seenHandler,
		"metadata":            _ObjTypeCoreV3EntityStateMetadataHandler,
		"sensu_agent_version": _ObjTypeCoreV3EntityStateSensu_agent_versionHandler,
		"system":              _ObjTypeCoreV3EntityStateSystemHandler,
	},
}