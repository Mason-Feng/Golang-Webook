local key= KEYS[1]
local cntKey=key..":cnt"

--你准备存储的验证码
local val=ARDV[1]

local ttl=tonumber(redis.call("ttl",key))

if ttl==-1 then

    --ttl=-1表示key存在，但是没有过期时间
    return -2

elseif ttl ==-2 or ttl<540 then
--可以发验证码
    redis.call("set",key,val)
    redis.call("expire",key,600)
    redis.call("set",cntKey,3)
    redis.call("set",cntKey,3)
else
    --发送太频繁
    return -1   
end   