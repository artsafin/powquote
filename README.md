## Running

Run a server:

    make run-server

Run clients:
    
    make run-client &
    make run-client &
    make run-client &

## PoW algorithm justification

Algorithm based on the http://hashcash.org/papers/dos-client-puzzles.pdf [1] paper which in turn is based on http://hashcash.org/papers/client-puzzles.pdf.

The algorithm used in this test task is Hashcash extended with a server interaction to retrieve a nonce and complexity, thus effectively making it a challenge-response protocol.

Why it was selected:
- (see main justification points in the paper [1], page 3)
- Hashcash is a CPU cost function. It was selected for the sake of the test task. Memory and network cost function may be more preferrable in the real world
- Hashcash is super easy for implementation on server and client side and it's implementation is easy to check
- complexity can be controlled by the server


## Known issues

For the sake of the test task simplicity:
- server does not change complexity dynamically depending on it's load 
- requests and responses are not signed
- client doesn't take into account server nonce timeout
- client doesn't retry if the solution is invalid (e.g. because of the previous point)

## Runtime configuration

### Server
`LISTEN` - interface and port to listen to (required)

`PROTECTED` - bool-ish value indicating DDoS protection enabled or not (default true)

`COMPLEXITY` - sets the static puzzle complexity (default 5) 

### Client

`SERVER` - address of the server (required)