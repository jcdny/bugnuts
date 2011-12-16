m2.xlarge:
  http://cloud.github.com/downloads/aichallenge/aichallenge/golang_60.1-9753~natty1_amd64.deb

engine.BenchmarkEngine	      10	 123288600 ns/op
engine.BenchmarkEngineOrdered	      10	 130909000 ns/op
game.BenchmarkParse	     100	  12438430 ns/op
maps.BenchmarkApplyCached	     200	  12557465 ns/op
maps.BenchmarkApplyNone	     100	  21649340 ns/op
maps.BenchmarkApplyNoCache	      50	  30984660 ns/op
maps.BenchmarkApplyCacheCreateA	      50	  34989800 ns/op
maps.BenchmarkApplyCacheCreateB	     100	  22765180 ns/op
maps.BenchmarkCacheAll	      10	 130144600 ns/op
maps.BenchmarkTransMap	    5000	    567735 ns/op
maps.BenchmarkTile0	      50	  31556800 ns/op
maps.BenchmarkTile4	      50	  30530820 ns/op
maps.BenchmarkTile8	      50	  30592700 ns/op
util.BenchmarkMinV	20000000	        79.6 ns/op
util.BenchmarkMin	20000000	        88.3 ns/op

bunty:

engine.BenchmarkEngine	       5	 206101200 ns/op
engine.BenchmarkEngineOrdered	       5	 212601600 ns/op
game.BenchmarkParse	     100	  14766500 ns/op
maps.BenchmarkApplyCached	     100	  15905120 ns/op
maps.BenchmarkApplyNone	     100	  27200050 ns/op
maps.BenchmarkApplyNoCache	      50	  52740880 ns/op
maps.BenchmarkApplyCacheCreateA	      50	  42257280 ns/op
maps.BenchmarkApplyCacheCreateB	      50	  28362680 ns/op
maps.BenchmarkCacheAll	      10	 152295300 ns/op
maps.BenchmarkTransMap	    2000	   1096505 ns/op
maps.BenchmarkTile0	      50	  34597160 ns/op
maps.BenchmarkTile4	      50	  34726680 ns/op
maps.BenchmarkTile8	      50	  34510780 ns/op
util.BenchmarkMinV	20000000	       101 ns/op
util.BenchmarkMin	20000000	        90.1 ns/op

Mac:

engine.BenchmarkEngine	      10	 124966000 ns/op
engine.BenchmarkEngineOrdered	      10	 131019800 ns/op
game.BenchmarkParse	     200	  12664680 ns/op
maps.BenchmarkApplyCached	     200	  11035745 ns/op
maps.BenchmarkApplyNone	     100	  22791480 ns/op
maps.BenchmarkApplyNoCache	      50	  28601600 ns/op
maps.BenchmarkApplyCacheCreateA	      50	  36333280 ns/op
maps.BenchmarkApplyCacheCreateB	     100	  25426170 ns/op
maps.BenchmarkCacheAll	      10	 144461100 ns/op
maps.BenchmarkTransMap	    5000	    654447 ns/op
maps.BenchmarkTile0	      50	  31541460 ns/op
maps.BenchmarkTile4	      50	  31571460 ns/op
maps.BenchmarkTile8	      50	  31571220 ns/op
util.BenchmarkMinV	50000000	        73.5 ns/op
util.BenchmarkMin	20000000	       107 ns/op

c1.xlarge:

game.BenchmarkParse	     100	  15677190 ns/op
maps.BenchmarkApplyCached	     100	  15758880 ns/op
maps.BenchmarkApplyNone	     100	  26980450 ns/op
maps.BenchmarkApplyNoCache	      50	  37851860 ns/op
maps.BenchmarkApplyCacheCreateA	      50	  42967940 ns/op
maps.BenchmarkApplyCacheCreateB	      50	  28350060 ns/op
maps.BenchmarkCacheAll	      10	 165462000 ns/op
maps.BenchmarkTransMap	    2000	    784545 ns/op
maps.BenchmarkTile0	      50	  39894780 ns/op
maps.BenchmarkTile4	      50	  39509760 ns/op
maps.BenchmarkTile8	      50	  39231520 ns/op
engine.BenchmarkEngine	      10	 156166400 ns/op
engine.BenchmarkEngineOrdered	      10	 168102500 ns/op
util.BenchmarkMinV	20000000	        99.1 ns/op
util.BenchmarkMin	20000000	       111 ns/op
