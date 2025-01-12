package postgres

const createOrUpdateEntityStateQuery = `
-- This query creates a new entity state, or updates it if it already exists.
--
-- The query has several WITH clauses, that allow it to affect multiple tables
-- in one shot. This allows us to use prepared statements, as they otherwise
-- won't work with multiple commands.
--
-- Parameters:
-- $1: The entity namespace.
-- $2: The entity name.
-- $3: The time the entity was last seen by the backend.
-- $4: The label selectors of the entity state.
-- $5: The annotations of the entity state.
-- $6-$17: System facts.
-- $18: The sensu agent version.
-- $19: Network names, an array of strings.
-- $20: Network MAC addresses, an array of strings.
-- $21: Network names, an array of json-encoded arrays.
-- $22: The entity state ID.
-- $23: The namespace ID.
-- $24: The entity config ID.
-- $25: The time that the entity state was created.
-- $26: The time that the entity state was last updated.
-- $27: The time that the entity state was soft deleted.
--
WITH ignored AS (
	SELECT
		$22::bigint,
		$25::timestamptz,
		$26::timestamptz,
		$27::timestamptz
), namespace AS (
	SELECT COALESCE (
		NULLIF($23, 0),
		(SELECT id FROM namespaces WHERE name = $1)
	) AS id
), config AS (
	SELECT COALESCE (
		NULLIF($24, 0),
		(SELECT id FROM entity_configs
			WHERE
				namespace_id = (SELECT id FROM namespace) AND
				name = $2
		)) AS id
), state AS (
	INSERT INTO entity_states (
		entity_config_id,
		namespace_id,
		name,
		expires_at,
		last_seen,
		selectors,
		annotations,
		deleted_at
	)
	VALUES (
		(SELECT id FROM config),
		(SELECT id FROM namespace),
		$2,
		now(),
		$3,
		$4,
		$5,
		NULL
	)
	ON CONFLICT ( namespace_id, name )
	DO UPDATE SET
		last_seen = $3,
		selectors = $4,
		annotations = $5,
		deleted_at = NULL
	RETURNING id AS id
), system AS (
	SELECT
	state.id AS entity_id,
	$6::text AS hostname,
	$7::text AS os,
	$8::text AS platform,
	$9::text AS platform_family,
	$10::text AS platform_version,
	$11::text AS arch,
	$12::integer AS arm_version, -- TODO(eric): is this a bug? Doesn't work without the cast.
	$13::text AS libc_type,
	$14::text AS vm_system,
	$15::text AS vm_role,
	$16::text AS cloud_provider,
	$17::text AS float_type,
	$18::text AS sensu_agent_version
	FROM state
), networks AS (
	SELECT
	nics.name AS name,
	nics.mac AS mac,
	nics.addresses AS addresses
	FROM UNNEST(
		cast($19 AS text[]),
		cast($20 AS text[]),
		cast($21 AS jsonb[])
	) AS nics ( name, mac, addresses )
),
sys_update AS (
	-- Insert or update the entity's system object.
	INSERT INTO entities_systems SELECT * FROM system
	ON CONFLICT (entity_id) DO UPDATE SET (
		hostname,
		os,
		platform,
		platform_family,
		platform_version,
		arch,
		arm_version,
		libc_type,
		vm_system,
		vm_role,
		cloud_provider,
		float_type,
		sensu_agent_version
	) = (
	-- Update the system if it exists already.
		SELECT
			hostname,
			os,
			platform,
			platform_family,
			platform_version,
			arch,
			arm_version,
			libc_type,
			vm_system,
			vm_role,
			cloud_provider,
			float_type,
			sensu_agent_version
		FROM system
	)
), del AS (
	-- Delete the networks that are not in the current set.
	DELETE FROM entities_networks
	USING state, networks
	WHERE entity_id = state.id
	AND ( entities_networks.name, entities_networks.mac, entities_networks.addresses ) NOT IN (SELECT * FROM networks)
)
-- Insert the networks that are new in the current set.
INSERT INTO entities_networks
SELECT state.id, networks.name, networks.mac, networks.addresses
FROM state, networks
ON CONFLICT ( entity_id, name, mac, addresses ) DO NOTHING
`

const createIfNotExistsEntityStateQuery = `
-- This query inserts rows into the entity_states table. By design, it
-- errors when an entity with the same namespace and name already
-- exists.
--
WITH ignored AS (
	SELECT
		$22::bigint,
		$25::timestamptz,
		$26::timestamptz,
		$27::timestamptz
), namespace AS (
	SELECT COALESCE (
		NULLIF($23, 0),
		(SELECT id FROM namespaces WHERE name = $1)
	) AS id
), config AS (
	SELECT COALESCE (
		NULLIF($24, 0),
		(SELECT id FROM entity_configs
			WHERE
				namespace_id = (SELECT id FROM namespace) AND
				name = $2
		)) AS id
), state AS (
	INSERT INTO entity_states (
		entity_config_id,
		namespace_id,
		name,
		expires_at,
		last_seen,
		selectors,
		annotations
	)
	VALUES (
		(SELECT id FROM config),
		(SELECT id FROM namespace),
		$2,
		now(),
		$3,
		$4,
		$5
)
	RETURNING id
), system AS (
	INSERT INTO entities_systems
	SELECT
	state.id,
	$6::text,
	$7::text,
	$8::text,
	$9::text,
	$10::text,
	$11::text,
	$12::integer, -- TODO(eric): is this a bug? Doesn't work without the cast.
	$13::text,
	$14::text,
	$15::text,
	$16::text,
	$17::text,
	$18::text
	FROM state
)
INSERT INTO entities_networks
SELECT
state.id,
nics.name,
nics.mac,
nics.addresses
FROM state, UNNEST(
	cast($19 AS text[]),
	cast($20 AS text[]),
	cast($21 AS jsonb[])
)
AS nics ( name, mac, addresses )
`

const updateIfExistsEntityStateQuery = `
-- This query updates the entity state, but only if it exists.
--
WITH ignored AS (
	SELECT
		$22::bigint,
		$23::bigint,
		$24::bigint,
		$25::timestamptz,
		$26::timestamptz,
		$27::timestamptz
), namespace AS (
	SELECT id FROM namespaces
	WHERE namespaces.name = $1
), state AS (
	SELECT id FROM entity_states
	WHERE
		entity_states.namespace_id = (SELECT id FROM namespace) AND
		entity_states.name = $2
), upd AS (
	UPDATE entity_states
	SET last_seen = $3, selectors = $4, annotations = $5, deleted_at = NULL
	FROM state
	WHERE state.id = entity_states.id
), system AS (
	SELECT
	state.id AS entity_id,
	$6::text AS hostname,
	$7::text AS os,
	$8::text AS platform,
	$9::text AS platform_family,
	$10::text AS platform_version,
	$11::text AS arch,
	$12::integer AS arm_version, -- TODO(eric): is this a bug? Doesn't work without the cast.
	$13::text AS libc_type,
	$14::text AS vm_system,
	$15::text AS vm_role,
	$16::text AS cloud_provider,
	$17::text AS float_type,
	$18::text AS sensu_agent_version
	FROM state
), networks AS (
	SELECT
	nics.name AS name,
	nics.mac AS mac,
	nics.addresses AS addresses
	FROM UNNEST(
		cast($19 AS text[]),
		cast($20 AS text[]),
		cast($21 AS jsonb[])
	)
	AS nics ( name, mac, addresses )
), sys_update AS (
	-- Insert or update the entity's system object.
	INSERT INTO entities_systems SELECT * FROM system
	ON CONFLICT (entity_id) DO UPDATE SET (
		hostname,
		os,
		platform,
		platform_family,
		platform_version,
		arch,
		arm_version,
		libc_type,
		vm_system,
		vm_role,
		cloud_provider,
		float_type,
		sensu_agent_version
	) = (
	-- Update the system if it exists already.
		SELECT
			hostname,
			os,
			platform,
			platform_family,
			platform_version,
			arch,
			arm_version,
			libc_type,
			vm_system,
			vm_role,
			cloud_provider,
			float_type,
			sensu_agent_version
		FROM system
	)
), del AS (
	-- Delete the networks that are not in the current set.
	DELETE FROM entities_networks
	USING state, networks
	WHERE entity_id = state.id
	AND ( entities_networks.name, entities_networks.mac, entities_networks.addresses ) NOT IN (SELECT * FROM networks)
), ins AS (
	-- Insert the networks that are new in the current set.
	INSERT INTO entities_networks
	SELECT state.id, networks.name, networks.mac, networks.addresses
	FROM state, networks
	ON CONFLICT ( entity_id, name, mac, addresses ) DO NOTHING
)
SELECT * FROM state;
`

const getEntityStateQuery = `
-- This query fetches a single entity state, or nothing.
--
WITH namespace AS (
	SELECT id FROM namespaces WHERE name = $1
), state AS (
	SELECT id FROM entity_states
	WHERE
		namespace_id = (SELECT id FROM namespace) AND
		name = $2 AND
		deleted_at IS NULL
), network AS (
	SELECT
		array_agg(entities_networks.name) AS names,
		array_agg(entities_networks.mac) AS macs,
		array_agg(entities_networks.addresses) AS addresses
	FROM entities_networks, state
	WHERE entities_networks.entity_id = state.id
)
SELECT
	$1,
	entity_states.name,
	entity_states.last_seen,
	entity_states.selectors,
	entity_states.annotations,
	entities_systems.hostname,
	entities_systems.os,
	entities_systems.platform,
	entities_systems.platform_family,
	entities_systems.platform_version,
	entities_systems.arch,
	entities_systems.arm_version,
	entities_systems.libc_type,
	entities_systems.vm_system,
	entities_systems.vm_role,
	entities_systems.cloud_provider,
	entities_systems.float_type,
	entities_systems.sensu_agent_version,
	network.names,
	network.macs,
	network.addresses::text[],
	entity_states.id,
	namespace.id,
	entity_states.entity_config_id,
	entity_states.created_at,
	entity_states.updated_at,
	entity_states.deleted_at
FROM entity_states LEFT OUTER JOIN entities_systems ON entity_states.id = entities_systems.entity_id, namespace, network, state
WHERE entity_states.id = state.id
`

const getEntityStatesQuery = `
-- This query fetches multiple entity states, including network information.
--
WITH namespace AS (
	SELECT id, name FROM namespaces
	WHERE namespaces.name = $1
), state AS (
	SELECT entity_states.id AS id
	FROM entity_states
	WHERE
		entity_states.namespace_id = (SELECT id FROM namespace) AND
		entity_states.name IN (SELECT unnest($2::text[])) AND
		entity_states.deleted_at IS NULL
), network AS (
	SELECT
		entities_networks.entity_id AS entity_id,
		array_agg(entities_networks.name) AS names,
		array_agg(entities_networks.mac) AS macs,
		array_agg(entities_networks.addresses) AS addresses
	FROM entities_networks, state
	WHERE entities_networks.entity_id = state.id
	GROUP BY entities_networks.entity_id
)
SELECT
	namespace.name,
	entity_states.name,
	entity_states.last_seen,
	entity_states.selectors,
	entity_states.annotations,
	entities_systems.hostname,
	entities_systems.os,
	entities_systems.platform,
	entities_systems.platform_family,
	entities_systems.platform_version,
	entities_systems.arch,
	entities_systems.arm_version,
	entities_systems.libc_type,
	entities_systems.vm_system,
	entities_systems.vm_role,
	entities_systems.cloud_provider,
	entities_systems.float_type,
	entities_systems.sensu_agent_version,
	network.names,
	network.macs,
	network.addresses::text[],
	entity_states.id,
	namespace.id,
	entity_states.entity_config_id,
	entity_states.created_at,
	entity_states.updated_at,
	entity_states.deleted_at
FROM entity_states
LEFT OUTER JOIN entities_systems ON entity_states.id = entities_systems.entity_id, namespace, network, state
WHERE entity_states.id = state.id AND entity_states.id = network.entity_id
`

const hardDeleteEntityStateQuery = `
-- This query deletes an entity state. Any related system & network state will
-- also be deleted via ON DELETE CASCADE triggers.
--
-- Parameters:
-- $1 Namespace
-- $2 Entity name
WITH namespace AS (
	SELECT id FROM namespaces
	WHERE name = $1
)
DELETE FROM entity_states
WHERE
	namespace_id = (SELECT id FROM namespace) AND
	name = $2;
`

const deleteEntityStateQuery = `
-- This query soft deletes an entity state.
--
-- Parameters:
-- $1 Namespace
-- $2 Entity name
WITH namespace AS (
	SELECT id FROM namespaces
	WHERE name = $1
)
UPDATE entity_states
SET deleted_at = now()
WHERE
	namespace_id = (SELECT id FROM namespace) AND
	name = $2;
`

const hardDeletedEntityStateQuery = `
-- This query discovers if an entity state has been hard deleted.
--
WITH namespace AS (
	SELECT id FROM namespaces
	WHERE name = $1
)
SELECT NOT EXISTS (
	SELECT true FROM entity_states
	WHERE
		namespace_id = (SELECT id FROM namespace) AND
		name = $2
);
`

const listEntityStateQuery = `
-- This query lists entity states from a given namespace.
--
WITH namespace AS (
	SELECT id, name FROM namespaces
	WHERE namespaces.name = $1 OR $1 IS NULL
), network AS (
	SELECT
		entity_states.id as entity_id,
		array_agg(entities_networks.name) AS names,
		array_agg(entities_networks.mac) AS macs,
		array_agg(entities_networks.addresses) AS addresses
	FROM entities_networks, entity_states
	WHERE entity_states.namespace_id = (SELECT id FROM namespace) OR $1 IS NULL
	GROUP BY entity_states.id
)
SELECT 
	namespace.name,
	entity_states.name,
	entity_states.last_seen,
	entity_states.selectors,
	entity_states.annotations,
	entities_systems.hostname,
	entities_systems.os,
	entities_systems.platform,
	entities_systems.platform_family,
	entities_systems.platform_version,
	entities_systems.arch,
	entities_systems.arm_version,
	entities_systems.libc_type,
	entities_systems.vm_system,
	entities_systems.vm_role,
	entities_systems.cloud_provider,
	entities_systems.float_type,
	entities_systems.sensu_agent_version,
	network.names,
	network.macs,
	network.addresses::text[],
	entity_states.id,
	namespace.id,
	entity_states.entity_config_id,
	entity_states.created_at,
	entity_states.updated_at,
	entity_states.deleted_at
FROM entity_states LEFT OUTER JOIN entities_systems ON entity_states.id = entities_systems.entity_id, network, namespace
WHERE
	network.entity_id = entity_states.id AND
	entity_states.deleted_at IS NULL
ORDER BY ( namespace.name, entity_states.name ) ASC
LIMIT $2
OFFSET $3
`

const listEntityStateDescQuery = `
-- This query lists entities from a given namespace.
--
WITH namespace AS (
	SELECT id, name FROM namespaces
	WHERE namespaces.name = $1 OR $1 IS NULL
), network AS (
	SELECT
		entity_states.id as entity_id,
		array_agg(entities_networks.name) AS names,
		array_agg(entities_networks.mac) AS macs,
		array_agg(entities_networks.addresses) AS addresses
	FROM entities_networks, entity_states
	WHERE entity_states.namespace_id = (SELECT id FROM namespace) OR $1 IS NULL
	GROUP BY entity_states.id
)
SELECT 
	namespace.name,
	entity_states.name,
	entity_states.last_seen,
	entity_states.selectors,
	entity_states.annotations,
	entities_systems.hostname,
	entities_systems.os,
	entities_systems.platform,
	entities_systems.platform_family,
	entities_systems.platform_version,
	entities_systems.arch,
	entities_systems.arm_version,
	entities_systems.libc_type,
	entities_systems.vm_system,
	entities_systems.vm_role,
	entities_systems.cloud_provider,
	entities_systems.float_type,
	entities_systems.sensu_agent_version,
	network.names,
	network.macs,
	network.addresses::text[],
	entity_states.id,
	namespace.id,
	entity_states.entity_config_id,
	entity_states.created_at,
	entity_states.updated_at,
	entity_states.deleted_at
FROM entity_states LEFT OUTER JOIN entities_systems ON entity_states.id = entities_systems.entity_id, network, namespace
WHERE
	network.entity_id = entity_states.id AND
	entity_states.deleted_at IS NULL
ORDER BY ( namespace.name, entity_states.name ) DESC
LIMIT $2
OFFSET $3
`

const existsEntityStateQuery = `
-- This query discovers if an entity exists, without retrieving it.
--
WITH namespace AS (
	SELECT id FROM namespaces
	WHERE name = $1
)
SELECT true FROM entity_states
WHERE
	namespace_id = (SELECT id FROM namespace) AND
	name = $2 AND
	deleted_at IS NULL;
`

const patchEntityStateQuery = `
-- This query updates only the last_seen and expires_at values of an
-- entity.
--
WITH namespace AS (
	SELECT id FROM namespaces
	WHERE name = $1
)
UPDATE entity_states
SET
	expires_at = $3 + now(),
	last_seen = $4,
	deleted_at = NULL
WHERE
	namespace_id = (SELECT id FROM namespace) AND
	name = $2;
`
