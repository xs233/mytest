#ifndef __INTERVAL_TREE_HPP__
#define __INTERVAL_TREE_HPP__

template <typename T1, typename T2>
class interval_tree {
public:
    interval_tree() {
        nil_ = new node{0};
        root_ = nil_;
    };

    ~interval_tree() {
        delete nil_;
        root_ = nullptr;
        nil_ = nullptr;
    };

    void insert(T1 low_endpoint, T1 high_endpoint, T2 value) {
        node* new_node = new node{low_endpoint, high_endpoint, value, RED, nil_, nil_};
        node* iterator_node = root_;
        node* last_iterator_node = nil_;
        while (iterator_node != nil_) {
            last_iterator_node = iterator_node;
            if (low_endpoint < iterator_node->low_endpoint) {
                iterator_node = iterator_node->left_node;
            } else {
                iterator_node = iterator_node->right_node;
            }
        }

        new_node->parent_node = last_iterator_node;
        if (last_iterator_node == nil_) {
            root_ = new_node;
        } else {
            if (low_endpoint < last_iterator_node->low_endpoint) {
                last_iterator_node->left_node = new_node;
            } else {
                last_iterator_node->right_node = new_node;
            }
        }
        insert_fixup(new_node);
    };

    void delete(node* useless_node) {
        if (root_ == nil_ || !useless_node || useless_node == nil_) {
            return;
        }

        node* replace_node = nullptr;
        COLOR original_color = useless_node->color;
        if (useless_node != nil_) {
            if (useless_node->left_node == nil_) {
                // No left node
                replace_node = useless_node->right_node;
                replace_useless_node(useless_node, replace_node);
            } else if (useless_node->right_node == nil_) {
                // No right node
                replace_node = useless_node->left_node;
                replace_useless_node(useless_node, replace_node);
            } else {
                // There are both left node and right node
                node* minimum_node = minimum(useless_node->right_node, nil_);
                replace_node = minimum_node->right_node;
                original_color = minimum_node->color;
                replace_useless_node(minimum_node, replace_node);
                useless_node->num = minimum_node->num;
                useless_node = minimum_node;
            }
        }

        if (original_color == BLACK) {
            delete_fixup(replace_node);
        }
    };

    T2& operator[](T1 low_endpoint, T1 high_endpoint) {
        if (low_endpoint == high_endpoint) {
            return nil_->value;
        }

        node* iterator_node = root_;
	    while (iterator_node != nil_) {
            if (low_endpoint == iterator_node->low_endpoint) {
                if (high_endpoint == iterator_node->high_endpoint) {
                    break;
                } else {
                    iterator_node = iterator_node->right_node;
                }
            } else if (num < iterator_node->num) {
                iterator_node = iterator_node->left_node;
            } else {
                iterator_node = iterator_node->right_node;
            }
        }
        return iterator_node->value;
    };

private:
    void left_rotate(node* x) {
        node* y = x->right_node;
        x->right_node = y->left_node;
        if (y->left_node != nil_) {
            y->left_node->parent_node = x;
        }

        y->parent_node = x->parent_node;
        if (x->parent_node == nil_) {
            root_ = y;
        } else if (x == x->parent_node->left_node) {
            x->parent_node->left_node = y;
        } else {
            x->parent_node->right_node = y;
        }
        y->left_node = x;
        x->parent_node = y;
    };

    void right_rotate(node* y) {
        node* x = y->left_node;
        y->left_node = x->right_node;
        if (x->right_node != nil_) {
            x->right_node->parent_node = y;
        }

        x->parent_node = y->parent_node;
        if (y->parent_node == nil_) {
            root_ = x;
        } else if (y == y->parent_node->left_node) {
            y->parent_node->left_node = x;
        } else {
            y->parent_node->right_node = x;
        }
        x->right_node = y;
        y->parent_node = x;
    };

    node* minimum(node* start, node* stop) {
        node* iterator_node = start;
        if (start == stop) {
            return nullptr;
        }

        while (iterator_node->left_node != stop) {
            iterator_node = iterator_node->left_node;
        }
        return iterator_node;
    };

    void insert_fixup(node* balance_node) {
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
                    left_rotate(balance_node);
                    balance_node->parent_node->color = BLACK;
                    balance_node->parent_node->parent_node->color = RED;
                    right_rotate(balance_node->parent_node->parent_node);
                } else {
                    balance_node->parent_node->color = BLACK;
                    balance_node->parent_node->parent_node->color = RED;
                    right_rotate(balance_node->parent_node->parent_node);
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
                    right_rotate(balance_node);
                    balance_node->parent_node->color = BLACK;
                    balance_node->parent_node->parent_node->color = RED;
                    left_rotate(balance_node->parent_node->parent_node);
                } else {
                    balance_node->parent_node->color = BLACK;
                    balance_node->parent_node->parent_node->color = RED;
                    left_rotate(balance_node->parent_node->parent_node);
                }
            }
        }
        root_->color = BLACK;
    };

    void replace_useless_node(node* useless_node, node* new_node) {
        if (useless_node->parent_node == nil_) {
            root_ = new_node;
        } else {
            if (useless_node == useless_node->parent_node->left_node) {
                useless_node->parent_node->left_node = new_node;
            } else {
                useless_node->parent_node->right_node = new_node;
            }
        }
        
        new_node->parent_node = useless_node->parent_node;
    };

    void rb_delete_fixup(node* balance_node) {
        while (balance_node != root_ && balance_node->color == BLACK) {
            if (balance_node == balance_node->parent_node->left_node) {
                node* brother_node = balance_node->parent_node->right_node;
                if (brother_node->color == RED) {
                    brother_node->color = BLACK;
                    brother_node->parent_node->color = RED;
                    left_rotate(brother_node->parent_node);
                    brother_node = balance_node->parent_node->right_node;
                }
                if (brother_node->left_node->color == BLACK && brother_node->right_node->color == BLACK) {
                    brother_node->color = RED;
                    balance_node = balance_node->parent_node;
                } else if (brother_node->right_node->color == BLACK) {
                    brother_node->left_node->color = BLACK;
                    brother_node->color = RED;
                    right_rotate(brother_node);
                    brother_node = balance_node->parent_node->right_node;
                } else {
                    brother_node->color = balance_node->parent_node->color;
                    balance_node->parent_node->color = BLACK;
                    brother_node->right_node->color = BLACK;
                    left_rotate(balance_node->parent_node);
                    balance_node = root_;
                }
            } else {
                node* brother_node = balance_node->parent_node->left_node;
                if (brother_node->color == RED) {
                    brother_node->color = BLACK;
                    brother_node->parent_node->color = RED;
                    right_rotate(brother_node->parent_node);
                    brother_node = balance_node->parent_node->left_node;
                }
                if (brother_node->left_node->color == BLACK && brother_node->right_node->color == BLACK) {
                    brother_node->color = RED;
                    balance_node = balance_node->parent_node;
                } else if (brother_node->left_node->color == BLACK) {
                    brother_node->right_node->color = BLACK;
                    brother_node->color = RED;
                    left_rotate(brother_node);
                    brother_node = balance_node->parent_node->left_node;
                } else {
                    brother_node->color = balance_node->parent_node->color;
                    balance_node->parent_node->color = BLACK;
                    brother_node->left_node->color = BLACK;
                    right_rotate(balance_node->parent_node);
                    balance_node = root_;
                }
            }
        }
        balance_node->color = BLACK;
    };

    enum COLOR {BLACK, RED};
    struct node {
        T1 low_endpoint;
        T1 high_endpoint;
        T2 value;
        COLOR color;
        node* left_node;
        node* right_node;
        node* parent_node;
    };

    node* root_;
    node* nil_;
};

#endif