go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof

Showing nodes accounting for -10512.39kB, 34.48% of 30492.55kB total
Dropped 3 nodes (cum <= 152.46kB)
flat  flat%   sum%        cum   cum%
-9928.45kB 32.56% 32.56% -12050.78kB 39.52%  compress/flate.NewWriter (inline)
-1642.35kB  5.39% 37.95% -1642.35kB  5.39%  compress/flate.(*compressor).initDeflate (inline)
544.67kB  1.79% 36.16% -2122.33kB  6.96%  compress/flate.(*compressor).init
513.12kB  1.68% 34.48%   513.12kB  1.68%  compress/flate.(*huffmanEncoder).generate
513kB  1.68% 32.79%      513kB  1.68%  runtime.allocm
-512.56kB  1.68% 34.48% -1024.65kB  3.36%  compress/flate.newHuffmanBitWriter (inline)
512.22kB  1.68% 32.80%   512.22kB  1.68%  runtime.malg
-512.09kB  1.68% 34.48%  -512.09kB  1.68%  compress/flate.newHuffmanEncoder (inline)
512.05kB  1.68% 32.80%   512.05kB  1.68%  context.(*cancelCtx).Done
512.05kB  1.68% 31.12%   512.05kB  1.68%  bufio.NewReaderSize (inline)
-512.03kB  1.68% 32.80%  -512.03kB  1.68%  syscall.anyToSockaddr
-512.02kB  1.68% 34.48%  -512.02kB  1.68%  time.map.init.0
0     0% 34.48%   512.05kB  1.68%  bufio.NewReader (inline)
0     0% 34.48%   513.12kB  1.68%  compress/flate.(*Writer).Close (inline)
0     0% 34.48%   513.12kB  1.68%  compress/flate.(*compressor).close
0     0% 34.48%   513.12kB  1.68%  compress/flate.(*compressor).storeHuff
0     0% 34.48%   513.12kB  1.68%  compress/flate.(*huffmanBitWriter).writeBlockHuff
0     0% 34.48%   513.12kB  1.68%  compress/gzip.(*Writer).Close
0     0% 34.48% -12050.78kB 39.52%  compress/gzip.(*Writer).Write
0     0% 34.48%   512.05kB  1.68%  database/sql.(*DB).connectionOpener
0     0% 34.48% -12050.78kB 39.52%  encoding/json.(*Encoder).Encode
0     0% 34.48% -12050.78kB 39.52%  github.com/alikhanturusbekov/go-url-shortener/internal/handler.(*URLHandler).ShortenURLAsJSON
0     0% 34.48% -26134.24kB 85.71%  github.com/alikhanturusbekov/go-url-shortener/pkg/compress.(*compressWriter).Write
0     0% 34.48% 14083.46kB 46.19%  github.com/alikhanturusbekov/go-url-shortener/pkg/compress.(*gzipResponseWriter).Write
0     0% 34.48% -11537.66kB 37.84%  github.com/go-chi/chi/v5.(*ChainHandler).ServeHTTP
0     0% 34.48% -11537.66kB 37.84%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
0     0% 34.48% -11537.66kB 37.84%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
0     0% 34.48% -12050.78kB 39.52%  github.com/go-chi/chi/v5/middleware.AllowContentType.func1.1
0     0% 34.48%  -512.03kB  1.68%  internal/poll.(*FD).Accept
0     0% 34.48%  -512.03kB  1.68%  internal/poll.accept
0     0% 34.48%  -512.03kB  1.68%  main.main
0     0% 34.48%  -512.03kB  1.68%  main.run
0     0% 34.48% -12050.78kB 39.52%  main.run.func1.AuthMiddleware.3.1
0     0% 34.48% -11537.66kB 37.84%  main.run.func1.GzipCompressor.2.1
0     0% 34.48% -11537.66kB 37.84%  main.run.func1.RequestLogger.1.1
0     0% 34.48%  -512.03kB  1.68%  net.(*TCPListener).Accept
0     0% 34.48%  -512.03kB  1.68%  net.(*TCPListener).accept
0     0% 34.48%  -512.03kB  1.68%  net.(*netFD).accept
0     0% 34.48%  -512.03kB  1.68%  net/http.(*Server).ListenAndServe
0     0% 34.48%  -512.03kB  1.68%  net/http.(*Server).Serve
0     0% 34.48% -11025.61kB 36.16%  net/http.(*conn).serve
0     0% 34.48% -11537.66kB 37.84%  net/http.HandlerFunc.ServeHTTP
0     0% 34.48%  -512.03kB  1.68%  net/http.ListenAndServe (inline)
0     0% 34.48%   512.05kB  1.68%  net/http.newBufioReader
0     0% 34.48% -11537.66kB 37.84%  net/http.serverHandler.ServeHTTP
0     0% 34.48%  -512.02kB  1.68%  runtime.doInit (inline)
0     0% 34.48%  -512.02kB  1.68%  runtime.doInit1
0     0% 34.48%     -513kB  1.68%  runtime.handoffp
0     0% 34.48% -1024.05kB  3.36%  runtime.main
0     0% 34.48%      513kB  1.68%  runtime.mcall
0     0% 34.48%      513kB  1.68%  runtime.newm
0     0% 34.48%   512.22kB  1.68%  runtime.newproc.func1
0     0% 34.48%   512.22kB  1.68%  runtime.newproc1
0     0% 34.48%      513kB  1.68%  runtime.park_m
0     0% 34.48%     1026kB  3.36%  runtime.resetspinning
0     0% 34.48%     -513kB  1.68%  runtime.retake
0     0% 34.48%     1026kB  3.36%  runtime.schedule
0     0% 34.48%      513kB  1.68%  runtime.startm
0     0% 34.48%     -513kB  1.68%  runtime.sysmon
0     0% 34.48%   512.22kB  1.68%  runtime.systemstack
0     0% 34.48%     1026kB  3.36%  runtime.wakep
0     0% 34.48%  -512.03kB  1.68%  syscall.Accept
0     0% 34.48%  -512.02kB  1.68%  time.init