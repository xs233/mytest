#include <cstdlib>
#include <ctime>
#include <iostream>

struct node {
	int num;
	node* left_node;
	node* right_node;
	node* parent;
};

void inorder_tree_walk(node* tree) {
	if (tree) {
		inorder_tree_walk(tree->left_node);
		std::cout << tree->num << '\n';
		inorder_tree_walk(tree->right_node);
	}
}

node* tree_search(node* tree, int num) {
	node* iterator_node = tree;
	while (iterator_node) {
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

node* tree_minimum(node* tree) {
	node* iterator_node = tree;
	if (!iterator_node) {
		return nullptr;
	}

	while (iterator_node->left_node) {
		iterator_node = iterator_node->left_node;
	}
	return iterator_node;
}

node* tree_maximum(node* tree) {
	node* iterator_node = tree;
	if (!iterator_node) {
		return nullptr;
	}

	while (iterator_node->right_node) {
		iterator_node = iterator_node->right_node;
	}
	return iterator_node;
}

void tree_insert(node** tree, node* new_node) {
	node* iterator_node = *tree;
	node* last_iterator_node = nullptr;
	while (iterator_node) {
		last_iterator_node = iterator_node;
		if (new_node->num < iterator_node->num) {
			iterator_node = iterator_node->left_node;
		}
		else {
			iterator_node = iterator_node->right_node;
		}
	}

	if (!last_iterator_node) {
		*tree = new_node;
	} else {
		if (new_node->num < last_iterator_node->num) {
			last_iterator_node->left_node = new_node;
		} else {
			last_iterator_node->right_node = new_node;
		}
	}
	new_node->parent = last_iterator_node;
	return;
}

void replace_useless_node(node** tree, node* useless_node, node* new_node) {
	if (useless_node->parent) {
		if (useless_node == useless_node->parent->left_node) {
			useless_node->parent->left_node = new_node;
		} else {
			useless_node->parent->right_node = new_node;
		}
	} else {
		*tree = new_node;
	}

	if (new_node) {
		new_node->parent = useless_node->parent;
	}
}

void tree_delete(node** tree, node* useless_node) {
	if (useless_node) {
		if (!useless_node->left_node) {
			// No left node
			replace_useless_node(tree, useless_node, useless_node->right_node);
		} else if (!useless_node->right_node) {
			// No right node
			replace_useless_node(tree, useless_node, useless_node->left_node);
		} else {
			// There are both left node and right node
			node* minimum_node = tree_minimum(useless_node->right_node);
			replace_useless_node(tree, minimum_node, minimum_node->right_node);
			useless_node->num = minimum_node->num;
			useless_node = minimum_node;
		}
		delete useless_node;
	}
}

int main() {
	node* tree = nullptr;
	std::srand(std::time(0));
	for (int i = 0; i < 50; ++i) {
		node* my_node = new node{ 0 };	
		my_node->num = std::rand() % 100;
		tree_insert(&tree, my_node);
	}

	std::cout << "print tree:" << '\n';
	inorder_tree_walk(tree);
	std::cout << "minimum value: " << tree_minimum(tree)->num << '\n';
	std::cout << "maximum value: " << tree_maximum(tree)->num << '\n';

	while (tree) {
		int i = std::rand() % 100;
		node* useless_node = tree_search(tree, i);
		if (useless_node) {
			std::cout << "Delete the the number of " << i << '\n';
			tree_delete(&tree, useless_node);
		}
		else {
			std::cout << "The tree don't contain the number of " << i << '\n';
		}
	}

	std::cout << "The tree is empty!" << '\n';
	return 0;
}