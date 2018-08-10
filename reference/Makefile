%.ll: %.c
	clang -O0 -S -emit-llvm -o $@ $^

clean:
	$(RM) *.ll
