#include <iostream>

int partition(int array[], int p, int r) {
	int q = p;
	for (int i = p; i < r - 1; ++i) {
		if (array[i] > array[r - 1]) {
			std::swap(array[i], array[q++]);
		}
	}
	std::swap(array[q], array[r - 1]);
	return q;
}

void quick_sort(int array[], int p, int r) {
	if (r > p + 1) {
		int q = partition(array, p, r);
		quick_sort(array, p, q);
		quick_sort(array, q + 1, r);
	}

}

int main() {
	int array[] = { 1,5,3,-9,-3,9,2,5,9,5,3,2,1 };
	quick_sort(array, 0, sizeof(array) / sizeof(int));
	for (int i = 0; i < sizeof(array) / sizeof(int); ++i) {
		std::cout << array[i] << '\n';
	}
	return 0;
}