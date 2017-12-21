#include <cstring>
#include <iostream>

struct range_values {
	float expectation;
	float probability;
	int root_num;
};

range_values** optimal_bst(float keywords[], int k_size, float pseudo_keywords[], int p_size) {
	if (p_size - k_size != 1) {
		return nullptr;
	}

	range_values** optimal_values = new range_values*[p_size];
	for (int i = 0; i < p_size; ++i) {
		optimal_values[i] = new range_values[p_size];
		std::memset(optimal_values[i], '\0', p_size * sizeof(range_values));
		optimal_values[i][i].expectation = pseudo_keywords[i];
		optimal_values[i][i].probability = pseudo_keywords[i];
	}

	for (int range_size = 1; range_size < p_size; ++range_size) {
		int max_left_endpont_location = p_size - range_size - 1;
		for (int left_endpoint_location = 0; left_endpoint_location <= max_left_endpont_location; ++left_endpoint_location) {
			int right_endpoint_location = range_size + left_endpoint_location;
			optimal_values[left_endpoint_location][right_endpoint_location].probability =
				optimal_values[left_endpoint_location][right_endpoint_location - 1].probability + keywords[right_endpoint_location - 1] + pseudo_keywords[right_endpoint_location];
			for (int index = left_endpoint_location; index < right_endpoint_location; ++index) {
				float temp_expectation = optimal_values[left_endpoint_location][index].expectation +
					optimal_values[index + 1][right_endpoint_location].expectation +
					optimal_values[left_endpoint_location][right_endpoint_location].probability;
				if (!optimal_values[left_endpoint_location][right_endpoint_location].expectation ||
					optimal_values[left_endpoint_location][right_endpoint_location].expectation > temp_expectation) {
					optimal_values[left_endpoint_location][right_endpoint_location].expectation = temp_expectation;
					optimal_values[left_endpoint_location][right_endpoint_location].root_num = index;
				}
			}
		}
	}
	return optimal_values;
}

void print_optimal_bst(range_values** optimal_values, int left_endpoint_location, int right_endpoint_location) {
	if (left_endpoint_location == right_endpoint_location) {
		return;
	}

	int optimal_num = optimal_values[left_endpoint_location][right_endpoint_location].root_num;
	std::cout << optimal_num << '\n';
	print_optimal_bst(optimal_values, left_endpoint_location, optimal_num);
	print_optimal_bst(optimal_values, optimal_num + 1, right_endpoint_location);
}

int main() {
	float keywords[] = { 0.15, 0.1, 0.05, 0.1, 0.2 };
	float pseudo_keywords[] = { 0.05, 0.1, 0.05, 0.05, 0.05, 0.1 };
	range_values** optimal_values = optimal_bst(keywords, sizeof(keywords) / sizeof(float), pseudo_keywords, sizeof(pseudo_keywords) / sizeof(float));
	print_optimal_bst(optimal_values, 0, sizeof(pseudo_keywords) / sizeof(float) - 1);
	for (int i = 0; i < sizeof(pseudo_keywords) / sizeof(float); ++i) {
		delete[] optimal_values[i];
	}
	delete[] optimal_values;
	return 0;
}