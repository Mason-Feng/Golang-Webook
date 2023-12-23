wrk.method="POST"
wrk.headers["Content-Type"]="application/json"
local random=math.random
local function uuid()
    local template='xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'
    return string.gsub( template,'[xy]',function(c)
        local v= ( c == 'x') and random(0,0xf) or random(8,0xb)
        return string.format('%x',v)
    end)
end

function init(args)
    cnt=0
    prefix=uuid()
end

function request()
    body=string.format('{"email":"%s%d@qq.com","password":"fqw#1234","confirmPassword":"fqw#1234"}',prefix,cnt)
    cnt=cnt+10
    return wrk.format('POST',wrk.path,wrk.headers,body)
end

function response()

end