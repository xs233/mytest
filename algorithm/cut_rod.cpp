#include <iostream>

int memoized_cut_rod_aux(int fixed_price[], int fixed_price_length, int length, int optimal_price[]) {
	if (!length || optimal_price[length]) {
		return optimal_price[length];
	}

	int max = 0;
	for (int i = 1; i <= (fixed_price_length < length ? fixed_price_length : length); ++i) {
		int temp = fixed_price[i - 1] + memoized_cut_rod_aux(fixed_price, fixed_price_length, length - i, optimal_price);
		max = max > temp ? max : temp;
	}
	optimal_price[length] = max;
	return max;
}

int memoized_cut_rod(int fixed_price[], int fixed_price_length, int length) {
	if (length <= 0) {
		return 0;
	}

	int* optimal_price = new int[length + 1] {0};
	int max = memoized_cut_rod_aux(fixed_price, fixed_price_length, length, optimal_price);
	delete[] optimal_price;
	return max;
}

int bottom_up_cut_rod(int fixed_price[], int fixed_price_length, int length) {
	if (length <= 0) {
		return 0;
	}

	int* optimal_price = new int[length + 1]{ 0 };
	int max = 0;
	for (int i = 1; i <= length; ++i) {
		for (int j = 1; j <= (fixed_price_length < i ? fixed_price_length : i); ++j) {
			int temp = fixed_price[j - 1] + optimal_price[i - j];
			max = max > temp ? max : temp;
		}
		optimal_price[i] = max;
	}

	delete[] optimal_price;
	return max;
}

int main() {
	int fixed_price[]{ 1,5,8,9,10,17,17,20,24,30 };
	std::cout << memoized_cut_rod(fixed_price, sizeof(fixed_price) / sizeof(int), 532) << '\n';
	std::cout << bottom_up_cut_rod(fixed_price, sizeof(fixed_price) / sizeof(int), 532) << '\n';
	return 0;
}