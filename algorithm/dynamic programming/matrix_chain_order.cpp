#include <iostream>
#include <array>
#include "interval_tree.hpp"

void matrix_chain_order(int matrix_chain[], int chain_length) {
    if (chain_length < 4) {
        return;
    }

    interval_tree<int, std::array<int, 2>> optimal_choice;
    for (int i = 2; i < chain_length; ++i) {
        for (int j = 1; j <= chain_length - i; ++j) {
            int k = i + j - 1;
            for (int l = j; l < k; ++l) {
                auto left_optimal_value = optimal_choice[j, l];
                auto right_optimal_value = optimal_choice[l + 1, k];
                auto optimal_value = optimal_choice[j, k];
                int temp =  left_optimal_value[0] + right_optimal_value[0] + matrix_chain[j - 1] * matrix_chain[l] * matrix_chain[k]; 
                if (optimal_value[0]) {
                    optimal_value[0] = temp < optimal_value[0] ? temp : optimal_value[0];
                    optimal_value[1] = l;
                } else {
                    optimal_choice.insert(j, k, std::array<int, 2>{temp, l});
                }
            }
        }
    }
}

int main() {
    int matrix_chain[] = {1, 100, 100, 100};
    matrix_chain_order(matrix_chain, sizeof(matrix_chain) / sizeof(int));
    return 0;
}