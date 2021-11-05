## VK-Tarantool

## Table of contents:
* [Run locally](#local)
* [Test deployed version](#deployed)
* [API endpoints](#endpoints)
* [Analytics](#analytics)

### How to run locally 💻:
<a name="local"></a>
1) Run tarantool in docker
```code
docker run --name mytarantool -d -p 3301:3301 \
tarantool/tarantool:2.6.0
```
⚠️ `if run on ARM - then use version 2.10.0-beta1`

2) Enter in container with docker exec
3) Create schemas:
```code
s = box.schema.space.create('short')

s:format({
{name = 'shorted', type = 'string'},
{name = 'original', type = 'string'},
{name = 'created_by', type = 'string'}
})

s:create_index('primary', {
         type = 'hash',
         parts = {'shorted'}
         })
```

```code
t = box.schema.space.create('transitions')

t:format({
{name = 'id', type = 'string'},
{name = 'shorted', type = 'string'},
{name = 'ip', type = 'string'},
{name = 'ua', type = 'string'}
})

t:create_index('primary', {
         type = 'hash',
         parts = {'id'}
         })
```

4) Create index
```code
t:create_index('shorted_idx', { type = 'tree', unique = false, parts = {'shorted'} })
```

5) Run with command:
```code
go run cmd/shorter/main.go
```
Server is now running on port 8182

### Deployed version:
<a name="deployed"></a>
Send http requests below to deployed version of service at address:
```http://37.139.34.190/```

### API Endpoints:
<a name="endpoints"></a>
1) Set endpoint
```code
/set?url=http://vk.com/ac
```
⚠️ URL must be with scheme (http:// or https://)

2) Get endpoint
```code
/{hash_from_previous_request}
```

### Analytics
<a name="analytics"></a>
Analytics is available from web UI after creating first short link