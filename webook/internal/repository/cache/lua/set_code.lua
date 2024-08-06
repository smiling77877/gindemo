---
--- Created by 95807.
--- DateTime: 2024/5/25 上午1:25
---
local key = KEYS[1]
local cntKey = key..":cnt"
-- 准备的存储的验证码
local val = ARGV[1]

local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
	-- key存在但是没有过期时间
	return -2
elseif ttl == -2 or ttl < 540 then
    -- 可以发验证码
    redis.call("set", key, val)
    -- 600秒
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    -- 发送太频繁
    return -1
end