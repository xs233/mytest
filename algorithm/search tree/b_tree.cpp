#include <iostream>
#include <tuple>

struct b_tree {
    int num;
    int* key;
    b_tree** value;
    bool leaf;
};

int t = 10;

auto b_tree_search(b_tree* root, int key) {
    for (int i = 0; i < root->num; ++i) {
        if (key == root->key[i]) {
            return std::make_tuple(root, i);
        } else if (key < root->key[i]) {
            return b_tree_search(root, key);
        }
    }
    return std::make_tuple(nullptr, 0);
}

void b_tree_split_child(b_tree* x, int m) {
    b_tree* y = x->value[m];
    b_tree* z = new b_tree{0};
    z->num = t - 1;
    z->leaf = y->leaf;
    for (int i = 0; i < z->num; ++i) {
        z->key[i] = y->key[t + i];
        if (!z->leaf) {
            z->value[i] = y->value[t + i];
        }
    }
    if (!z->leaf) {
        z->value[t - 1] = y->value[2 * t - 1];
    }
    y->num = t - 1;

    ++x->num;
    for (int i = x->num - 1; i > m; --i) {
        x->key[i] = x->key[i - 1];
        x->value[i + 1] = x->value[i];
    }
    x->key[m] = y->key[t - 1];
    x->value[m] = y;
    x->value[m + 1] = z;
}

int main() {

    return 0;
}