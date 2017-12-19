#include <cstring>
#include <iostream>


void counting_sort(int original_array[], int sort_array[], int array_size, int min, int max) {
	int range_size = max - min + 1;
	int* counting_array = new int[range_size];
	std::memset(counting_array, 0, range_size * sizeof(int));

	// Counting
	for (unsigned int i = 0; i < array_size; ++i) {
		++counting_array[original_array[i] - min];
	}

	// Records the number of elements before each location
	for (unsigned int i = 1; i < range_size; ++i) {
		counting_array[i] += counting_array[i - 1];
	}

	// Start sorting
	for (int i = array_size - 1; i >= 0; --i) {
		sort_array[--counting_array[original_array[i] - min]] = original_array[i];
	}
	delete[] counting_array;
}

int main() {
	int original_array[] = { -1,-3,-5,8,3,2,-9,6,0,2,4,1 };
	int array_size = sizeof(original_array) / sizeof(int);
	int* sort_array = new int[array_size];
	counting_sort(original_array, sort_array, array_size, -9, 8);

	for (unsigned int i = 0; i < array_size; ++i) {
		std::cout << sort_array[i] << '\n';
	}
	delete[] sort_array;
	return 0;
}