blake3:
	gcc -shared -O3 -o libblake3.so -DBLAKE3_USE_NEON blake3.c blake3_dispatch.c blake3_portable.c blake3_neon.c
