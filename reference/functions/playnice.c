// functions aren't closures so they don't step on each others variables
int foo() {
    int x = 2;
    return x;
}
int main() {
    int x = 1;
    int y = foo();
    return x;
}
