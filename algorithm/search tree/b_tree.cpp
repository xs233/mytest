#include <iostream>
#include <tuple>

struct b_tree {
    int num;
    int* key;
    b_tree* value;
    bool leaf;
};

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

int main() {

    return 0;
}