// File scripts.go contains code related to parsing
// lua scripts in the scripts file.

// This file has been automatically generated by go generate,
// which calls scripts/main.go. Do not edit it directly!

package kvmodel

import (
	"github.com/garyburd/redigo/redis"
)

var (
	deleteModelsBySetIdsScript = redis.NewScript(0, `-- Copyright 2015 Alex Browne.  All rights reserved.
-- Use of this source code is governed by the MIT
-- license, which can be found in the LICENSE file.

-- delete_models_by_set_ids is a lua script that takes the following arguments:
-- 	1) The key of a set of model ids
--		2) The name of a registered model
-- The script then deletes all the models corresponding to the ids in the given
-- set. It returns the number of models that were deleted. It does not delete the
-- given set.

-- IMPORTANT: If you edit this file, you must run go generate . to rewrite ../scripts.go

-- Assign keys to variables for easy access
local setKey = ARGV[1]
local collectionName = ARGV[2]
-- Get all the ids from the set name
local ids = redis.call('SMEMBERS', setKey)
local count = 0
if #ids > 0 then
	-- Iterate over the ids
	for i, id in ipairs(ids) do
		-- Delete the main hash for each model
		local key = collectionName .. ':' .. id
		count = count + redis.call('DEL', key)
		-- Remove the model id from the set of all ids
		-- NOTE: this is not necessarily the same as the
		-- setName we were given
		local setKey = collectionName .. ':all'
		redis.call('SREM', setKey, id)
	end
end
return count
`)
	deleteStringIndexScript = redis.NewScript(0, `-- Copyright 2015 Alex Browne.  All rights reserved.
-- Use of this source code is governed by the MIT
-- license, which can be found in the LICENSE file.

-- delete_string_index is a lua script that takes the following arguments:
-- 	1) The name of a registered model
--		2) The id of the model to be deleted from the index
--		3) The name of the indexed string field
-- The script then checks if there is a value for the given field name stored in the
-- model hash, and if there is, removes the model from the index on the given field.
-- NOTE: This script *must* be called before the main hash for the model is updated/deleted.

-- IMPORTANT: If you edit this file, you must run go generate . to rewrite ../scripts.go

-- Assign keys to variables for easy access
local collectionName = ARGV[1]
local modelID = ARGV[2]
local fieldName = ARGV[3]
-- Get the old value from the existing model hash (if any)
local modelKey = collectionName .. ":" .. modelID
local oldValue = redis.call("HGET", modelKey, fieldName)
local indexKey = collectionName .. ":" .. fieldName
if oldValue ~= false then
	-- Remove the model from the field index
	local oldMember = oldValue .. "\0" .. modelID
	redis.call("ZREM", indexKey, oldMember)
end
`)
	extractIdsFromFieldIndexScript = redis.NewScript(0, `-- Copyright 2015 Alex Browne.  All rights reserved.
-- Use of this source code is governed by the MIT
-- license, which can be found in the LICENSE file.

-- exctract_ids_from_field_index is a lua script that takes the following arguments:
-- 	1) setKey: The key of a sorted set for a field index (either numeric or bool)
-- 	2) destKey: The key of a sorted set where the resulting ids will be stored
--		3) min: The min argument for the ZRANGEBYSCORE command
-- 	4) max: The max argument for the ZRANGEBYSCORE command
-- The script then calls ZRANGEBYSCORE on setKey with the given min and max arguments,
-- and then stores the resulting set in destKey. It does not preserve the existing
-- scores, and instead just replaces scores with sequential numbers to keep the members
-- in the same order.

-- IMPORTANT: If you edit this file, you must run go generate . to rewrite ../scripts.go

-- Assign keys to variables for easy access
local setKey = ARGV[1]
local destKey = ARGV[2]
local min = ARGV[3]
local max = ARGV[4]
-- Get all the members (value+id pairs) from the sorted set
local members = redis.call('ZRANGEBYSCORE', setKey, min, max)
-- Iterate over the members and add each to the destKey
for i, member in ipairs(members) do
	redis.call('ZADD', destKey, i, member)
end
`)
	extractIdsFromStringIndexScript = redis.NewScript(0, `-- Copyright 2015 Alex Browne.  All rights reserved.
-- Use of this source code is governed by the MIT
-- license, which can be found in the LICENSE file.

-- exctract_ids_from_string_index is a lua script that takes the following arguments:
-- 	1) setKey: The key of a sorted set for a string index, where each member is of the
--			form: value + NULL + id, where NULL is the ASCII NULL character which has a codepoint
--			value of 0.
--		2) destKey: The key of a sorted set where the resulting ids will be stored
-- 	3) min: The min argument for the ZRANGEBYLEX command
-- 	4) max: The max argument for the ZRANGEBYLEX command
-- The script then extracts the ids from setKey using the given min and max arguments,
-- and then stores them destKey with the appropriate scores in ascending order.

-- IMPORTANT: If you edit this file, you must run go generate . to rewrite ../scripts.go

-- Assign keys to variables for easy access
local setKey = ARGV[1]
local destKey = ARGV[2]
local min = ARGV[3]
local max = ARGV[4]
-- Get all the members (value+id pairs) from the sorted set
local members = redis.call('ZRANGEBYLEX', setKey, min, max)
if #members > 0 then
	-- Iterate over the members and extract the ids
	for i, member in ipairs(members) do
		-- The id is everything after the last space
		-- Find the index of the last space
		local idStart = string.find(member, '%z[^%z]*$')
		local id = string.sub(member, idStart+1)
		redis.call('ZADD', destKey, i, id)
	end
end
`)
)
