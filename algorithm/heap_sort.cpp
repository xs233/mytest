#include <iostream>

void max_heapify(int array[], int i, int size) {
	int left = i * 2 + 1;
	int right = i * 2 + 2;
	// No child node
	if (size <= left) {
		return;
	}

	// Just has left node
	if (size == right) {
		if (array[i] < array[left]) {
			std::swap(array[i], array[left]);
		}
		return;
	}

	// There are both left node and right node.
	if (array[left] > array[i] && array[left] >= array[right]) {
		std::swap(array[left], array[i]);
		max_heapify(array, left, size);
	}
	else if (array[right] > array[i] && array[right] > array[left]) {
		std::swap(array[right], array[i]);
		max_heapify(array, right, size);
	}

	return;
}

void build_max_heap(int array[], int size) {
	for (int i = size - 1; i > 0; --i) {
		max_heapify(array, (i - 1) / 2, size);
	}
}

void heap_sort(int array[], int size) {
	build_max_heap(array, size);
	for (int i = size - 1; i > 0; --i) {
		std::swap(array[0], array[i]);
		max_heapify(array, 0, i);
	}
}

int main() {
	int array[] = { 1, 2, 3, 10, 9, 3, 5, -2, -5, 11 };
	int size = sizeof(array) / sizeof(int);
	heap_sort(array, size);

	for (int i = 0; i < size; ++i) {
		std::cout << array[i] << " ";
	}
	return 0;
}