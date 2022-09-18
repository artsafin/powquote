## Running

Run a server:

    make run-server

Run clients:
    
    make run-client &
    make run-client &
    make run-client &

## PoW algorithm justification

Algorithm based on the http://hashcash.org/papers/dos-client-puzzles.pdf [1] paper which in turn is based on http://hashcash.org/papers/client-puzzles.pdf.

The algorithm used in this test task is the Hashcash extended with a server interaction to retrieve a nonce and complexity, thus effectively making it conform to a challenge-response protocol.

Why it was selected:
- (see main justification points in the paper [1], page 3)
- Hashcash is a CPU cost function. It was selected for the sake of the test task. Memory and network cost function may be more preferrable in the real world
- Hashcash is super easy for implementation on server and client side and it's implementation is easy to check
- complexity can be controlled by the server

###### Brief description of the algorithm logic and flow

1. Client knows server address; It makes `hello` request to the server
2. Server generates a `server nonce` and `complexity` value depending on the current load (note: complexity is static in current implementation) and sends both values back to client;
    
    Then server forcefully closes TCP connection with the client to save resources (because we are deterring a DoS attack, aren't we?)
    - `server nonce` is a uint64 number;
    - `complexity` is an int in range of `[0; 40]` where `0` complexity means "protection disabled", and `40` complexity means "impossible to solve".
3. Client generates it's own `client nonce` and starts a process of puzzle solving.
4. Client puzzle solving process is generating a pack of `random bytes` so that a `sha1(client IP, server nonce, client nonce, random bytes)` represented as hex string will turn out to contain sequential leading zero characters (literally `"0"`, not `"\0"`).
    The number of required zero characters to fulfill the puzzle is a `complexity` provided by the server initially.
5. Once puzzle is solved all the inputs of the hash function are sent to the server to check it's validity
6. Solving puzzle requires a CPU work on the client side
7. Server runs the same hash function and checks the number of zero leading characters. If it agrees that the complexity was met it provides access to it's resources
8. `server nonce` is changed every 5 minutes, effectively giving client only 5 minutes to solve the puzzle

## Known issues

For the sake of the test task simplicity:
- server does not change complexity dynamically depending on it's load 
- requests and responses are not signed
- client doesn't take into account server nonce timeout
- client doesn't retry if the solution is invalid (e.g. because of the previous point or server was restarted)

## Runtime configuration

### Server
`LISTEN` - interface and port to listen to (required)

`PROTECTED` - bool-ish value indicating DDoS protection enabled or not (default true)

`COMPLEXITY` - sets the static puzzle complexity (default 5) 

### Client

`SERVER` - address of the server (required)

`VERBOSE` - address of the server (required)