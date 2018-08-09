%.ll: %.c
	clang -S -emit-llvm -o $@ $^

clean:
	$(RM) *.ll
