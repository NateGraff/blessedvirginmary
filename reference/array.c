int main() {
    int x[5];
    x[0] = 1;
    x[1] = 2;
    x[2] = 3;
    x[3] = 4;
    x[4] = 5;
    int y = x[0];
    y += x[1];
    y += x[2];
    y += x[3];
    y += x[4];
    return y;
}
