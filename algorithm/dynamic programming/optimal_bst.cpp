#include <cstring>
#include <iostream>

struct range_values {
    float expectation;
    float probability;
};

void optimal_bst(float keywords[], int k_size, float pseudo_keywords[], int p_size) {
    if (p_size - k_size != 1) {
        return;
    }

    range_values** optimal_values = new range_values*[p_size];
    for (int i = 0; i < p_size; ++i) {
        optimal_values[i] = new range_values[k_size];
        std::memset(optimal_values[i], '\0', k_size * sizeof(range_values));
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
                if (left_endpoint_location > index - 1) {
                    optimal_values[left_endpoint_location][right_endpoint_location].expectation += 
                        optimal_values[left_endpoint_location][index].expectation;
                } else {
                    optimal_values[left_endpoint_location][right_endpoint_location].expectation += 
                        optimal_values[left_endpoint_location][index - 1].expectation;
                }
                optimal_values[left_endpoint_location][right_endpoint_location].expectation = 
                    optimal_values[index + 1][right_endpoint_location].expectation + 
                    optimal_values[left_endpoint_location][right_endpoint_location].probability;
            }
            
        }
    }
}

int main() {


    return 0;
}