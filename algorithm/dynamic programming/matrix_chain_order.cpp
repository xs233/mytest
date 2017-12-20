#include <iostream>
#include <array>
#include "interval_tree.hpp"

void matrix_chain_order(int matrix_chain[], int chain_length, interval_tree<int, std::array<int, 2>>& optimal_choice) {
	if (chain_length < 4) {
		return;
	}

	for (int i = 2; i < chain_length; ++i) {
		for (int j = 1; j <= chain_length - i; ++j) {
			int k = i + j - 1;
			for (int l = j; l < k; ++l) {
				auto& left_optimal_value = optimal_choice.get(j, l);
				auto& right_optimal_value = optimal_choice.get(l + 1, k);
				auto& optimal_value = optimal_choice.get(j, k);
				int temp = left_optimal_value[0] + right_optimal_value[0] + matrix_chain[j - 1] * matrix_chain[l] * matrix_chain[k];
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

void print_optimal_parens(interval_tree<int, std::array<int, 2>>& optimal_choice, int low_endpoint, int high_endpoint) {
	if (low_endpoint == high_endpoint) {
		std::cout << "A" << low_endpoint;
	} else {
		std::cout << "(";
		print_optimal_parens(optimal_choice, low_endpoint, optimal_choice.get(low_endpoint, high_endpoint)[1]);
		print_optimal_parens(optimal_choice, optimal_choice.get(low_endpoint, high_endpoint)[1] + 1, high_endpoint);
		std::cout << ")";
	}                    
}                

int main() { 
	int matrix_chain[] = { 1, 100, 100, 100 };
	interval_tree<int, std::array<int, 2>> optimal_choice;
	matrix_chain_order(matrix_chain, sizeof(matrix_chain) / sizeof(int), optimal_choice);
	print_optimal_parens(optimal_choice, 1, sizeof(matrix_chain) / sizeof(int) - 1);
	return 0;
}