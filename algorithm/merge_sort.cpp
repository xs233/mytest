#include <cstring>
#include <iostream>
using namespace std;

void merge(int* array, int p, int q, int r) {
	int n1 = q - p;
	int n2 = r - q;
	int* l_array = new int[n1];
	int* r_array = new int[n2];
	memcpy(l_array, array + p, n1 * sizeof(int));
	memcpy(r_array, array + q, n2 * sizeof(int));

	int i = 0; 
	int j = 0;
	for (int k = p; k < r; ++k) {
		if (i == n1) {
			array[k] = r_array[j++];
		}
		else if (j == n2) {
			array[k] = l_array[i++];
		} else {
			if (l_array[i] <= r_array[j]) {
				array[k] = l_array[i++];
			} else {
				array[k] = r_array[j++];
			}
		}
	}

	delete[] l_array;
	delete[] r_array;
}

void merge_sort(int* array, int p, int r) {
	if (r - p > 1) {
		int q = (p + r) / 2;
		merge_sort(array, p, q);
		merge_sort(array, q, r);
		merge(array, p, q, r);
	}
}

int main() {
	int array[] = { 3,2,6,5,1,2,2,2,0,-1,8,3,4,2,1,9,10,5,5,5,3,3,2,2,1,1,4,4,99,100 };
	merge_sort(array, 0, sizeof(array) / sizeof(int));
	for (int i = 0; i < sizeof(array) / sizeof(int); ++i) {
		cout << array[i] << '\n';
	}
	return 0;
}