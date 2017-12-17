#include <cstdlib>
#include <ctime>
#include <iostream>

enum COLOR {RED, BLACK};

struct node {
	int num;
    COLOR color;
	node* left_node;
	node* right_node;
	node* parent_node;
};

struct red_black_tree {
    node* root;
    node* nil;
};

void left_rotate(red_black_tree* tree, node* x) {
    node* y = x->right_node;
    x->right_node = y->left_node;
    if (y->left_node != tree->nil) {
        y->left_node->parent_node = x;
    }

    y->parent_node = x->parent_node;
    if (x->parent_node == tree->nil) {
        tree->root = y;
    } else if (x == x->parent_node->left_node) {
        x->parent_node->left_node = y;
    } else {
        x->parent_node->right_node = y;
    }
    y->left_node = x;
    x->parent_node = y;
}

void right_rotate(red_black_tree* tree, node* y) {
    node* x = y->left_node;
    y->left_node = x->right_node;
    if (x->right_node != tree->nil) {
        x->right_node->parent_node = y;
    }

    x->parent_node = y->parent_node;
    if (y->parent_node == tree->nil) {
        tree->root = x;
    } else if (y == y->parent_node->left_node) {
        y->parent_node->left_node = x;
    } else {
        y->parent_node->right_node = x;
    }
    x->right_node = y;
    y->parent_node = x;
}

void inorder_tree_walk(node* start, node* stop) {
	if (start != stop) {
		inorder_tree_walk(start->left_node, stop);
		std::cout << start->num << '\n';
		inorder_tree_walk(start->right_node, stop);
	}
}

node* rb_minimum(node* start, node* stop) {
	node* iterator_node = start;
	if (start == stop) {
		return nullptr;
	}

	while (iterator_node->left_node != stop) {
		iterator_node = iterator_node->left_node;
	}
	return iterator_node;
}

node* rb_maximum(node* start, node* stop) {
	node* iterator_node = start;
	if (start == stop) {
		return nullptr;
	}

	while (iterator_node->right_node != stop) {
		iterator_node = iterator_node->right_node;
	}
	return iterator_node;
}

void rb_insert_fixup(red_black_tree* tree, node* balance_node) {
    while (balance_node->parent_node->color == RED) {
        if (balance_node->parent_node == balance_node->parent_node->parent_node->left_node) {
            node* uncle_node = balance_node->parent_node->parent_node->right_node;
            if (uncle_node->color == RED) {
                balance_node->parent_node->parent_node->color = RED;
                balance_node->parent_node->color = BLACK;
                uncle_node->color = BLACK;
                balance_node = balance_node->parent_node->parent_node;
            } else if (balance_node == balance_node->parent_node->right_node) {
                balance_node = balance_node->parent_node;
                left_rotate(tree, balance_node);
                balance_node->parent_node->color = BLACK;
                balance_node->parent_node->parent_node->color = RED;
                right_rotate(tree, balance_node->parent_node->parent_node);
            } else {
                balance_node->parent_node->color = BLACK;
                balance_node->parent_node->parent_node->color = RED;
                right_rotate(tree, balance_node->parent_node->parent_node);
            }
        } else {
            node* uncle_node = balance_node->parent_node->parent_node->left_node;
            if (uncle_node->color == RED) {
                balance_node->parent_node->parent_node->color = RED;
                balance_node->parent_node->color = BLACK;
                uncle_node->color = BLACK;
                balance_node = balance_node->parent_node->parent_node;
            } else if (balance_node == balance_node->parent_node->left_node) {
                balance_node = balance_node->parent_node;
                right_rotate(tree, balance_node);
                balance_node->parent_node->color = BLACK;
                balance_node->parent_node->parent_node->color = RED;
                left_rotate(tree, balance_node->parent_node->parent_node);
            } else {
                balance_node->parent_node->color = BLACK;
                balance_node->parent_node->parent_node->color = RED;
                left_rotate(tree, balance_node->parent_node->parent_node);
            }
        }
    }
    tree->root->color = BLACK;
}

node* rb_search(red_black_tree* tree, int num) {
	node* iterator_node = tree->root;
	while (iterator_node != tree->nil) {
		if (num == iterator_node->num) {
			break;
		} else if (num > iterator_node->num) {
			iterator_node = iterator_node->right_node;
		} else {
			iterator_node = iterator_node->left_node;
		}
	}
	return iterator_node;
}

void rb_insert(red_black_tree* tree, node* new_node) {
    node* iterator_node = tree->root;
    node* last_iterator_node = tree->nil;
    while (iterator_node != tree->nil) {
        last_iterator_node = iterator_node;
        if (new_node->num < iterator_node->num) {
            iterator_node = iterator_node->left_node;
        } else {
            iterator_node = iterator_node->right_node;
        }
    }

    new_node->parent_node = last_iterator_node;
    if (last_iterator_node == tree->nil) {
        tree->root = new_node;
    } else {
        if (new_node->num < last_iterator_node->num) {
            last_iterator_node->left_node = new_node;
        } else {
            last_iterator_node->right_node = new_node;
        }
    }
    new_node->color = RED;
    new_node->left_node = tree->nil;
    new_node->right_node = tree->nil;
    rb_insert_fixup(tree, new_node);
}

void replace_useless_node(red_black_tree* tree, node* useless_node, node* new_node) {
	if (useless_node->parent_node == tree->nil) {
		tree->root = new_node;
	} else {
		if (useless_node == useless_node->parent_node->left_node) {
			useless_node->parent_node->left_node = new_node;
		} else {
			useless_node->parent_node->right_node = new_node;
		}
	}
	
	new_node->parent_node = useless_node->parent_node;
}

void rb_delete_fixup(red_black_tree* tree, node* balance_node) {
    while (balance_node != tree->root && balance_node->color == BLACK) {
        if (balance_node == balance_node->parent_node->left_node) {
            node* brother_node = balance_node->parent_node->right_node;
            if (brother_node->color == RED) {
                brother_node->color = BLACK;
                brother_node->parent_node->color = RED;
                left_rotate(tree, brother_node->parent_node);
                brother_node = balance_node->parent_node->right_node;
            }
            if (brother_node->left_node->color == BLACK && brother_node->right_node->color == BLACK) {
                brother_node->color = RED;
                balance_node = balance_node->parent_node;
            } else if (brother_node->right_node->color == BLACK) {
                brother_node->left_node->color = BLACK;
                brother_node->color = RED;
                right_rotate(tree, brother_node);
                brother_node = balance_node->parent_node->right_node;
            } else {
                brother_node->color = balance_node->parent_node->color;
                balance_node->parent_node->color = BLACK;
                brother_node->right_node->color = BLACK;
                left_rotate(tree, balance_node->parent_node);
                balance_node = tree->root;
            }
        } else {
            node* brother_node = balance_node->parent_node->left_node;
            if (brother_node->color == RED) {
                brother_node->color = BLACK;
                brother_node->parent_node->color = RED;
                right_rotate(tree, brother_node->parent_node);
                brother_node = balance_node->parent_node->left_node;
            }
            if (brother_node->left_node->color == BLACK && brother_node->right_node->color == BLACK) {
                brother_node->color = RED;
                balance_node = balance_node->parent_node;
            } else if (brother_node->left_node->color == BLACK) {
                brother_node->right_node->color = BLACK;
                brother_node->color = RED;
                left_rotate(tree, brother_node);
                brother_node = balance_node->parent_node->left_node;
            } else {
                brother_node->color = balance_node->parent_node->color;
                balance_node->parent_node->color = BLACK;
                brother_node->left_node->color = BLACK;
                right_rotate(tree, balance_node->parent_node);
                balance_node = tree->root;
            }
        }
    }
    balance_node->color = BLACK;
}

void rb_delete(red_black_tree* tree, node* useless_node) {
    if (!tree || tree->root == tree->nil || !useless_node || useless_node == tree->nil) {
        return;
    }

	node* replace_node = nullptr;
	COLOR original_color = useless_node->color;
	if (useless_node != tree->nil) {
		if (useless_node->left_node == tree->nil) {
			// No left node
			replace_node = useless_node->right_node;
			replace_useless_node(tree, useless_node, replace_node);
		} else if (useless_node->right_node == tree->nil) {
			// No right node
			replace_node = useless_node->left_node;
			replace_useless_node(tree, useless_node, replace_node);
		} else {
			// There are both left node and right node
			node* minimum_node = rb_minimum(useless_node->right_node, tree->nil);
			replace_node = minimum_node->right_node;
			original_color = minimum_node->color;
			replace_useless_node(tree, minimum_node, replace_node);
			useless_node->num = minimum_node->num;
			useless_node = minimum_node;
		}
	}

    if (original_color == BLACK) {
        rb_delete_fixup(tree, replace_node);
    }
}

int main() {
    node sentinel_node{0};
    sentinel_node.color = BLACK;
    red_black_tree* tree = new red_black_tree{0};
    tree->nil = &sentinel_node;
    tree->root = tree->nil;
    std::srand(std::time(0));
    for (int i = 0; i < 20; ++i) {
        node* new_node = new node{0};
        new_node->num = std::rand() % 100;
        rb_insert(tree, new_node);
    }
    inorder_tree_walk(tree->root, tree->nil);
    while (tree->root != tree->nil) {
        int i = std::rand() % 100;
		node* useless_node = rb_search(tree, i);
		if (useless_node != tree->nil) {
			std::cout << "Delete the the number of " << i << '\n';
			rb_delete(tree, useless_node);
		} else {
			std::cout << "The tree don't contain the number of " << i << '\n';
		}
    }
    return 0;
}