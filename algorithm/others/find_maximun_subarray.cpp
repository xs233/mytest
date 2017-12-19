#include <iostream>
#include <tuple>

std::tuple<int, int, int> find_max_crossing_subarray(int array[], int low, int mid, int high) {
	int sum = 0;
	int left_sum = array[mid - 1];
	int max_left = mid - 1;
	for (int i = mid - 1; i >= low; --i) {
		sum += array[i];
		if (sum > left_sum) {
			left_sum = sum;
			max_left = i;
		}
	}

	sum = 0;
	int right_sum = array[mid];
	int max_right = mid + 1;
	for (int i = mid; i < high; ++i) {
		sum += array[i];
		if (sum > right_sum) {
			right_sum = sum;
			max_right = i + 1;
		}
	}

	return std::make_tuple(max_left, max_right, left_sum + right_sum);
}

std::tuple<int, int, int> find_maximun_subarray(int array[], int low, int high) {
	if (high - low == 1) {
		return std::make_tuple(low, high, array[low]);
	}

	int mid = (low + high) / 2;
	auto left_maximun_subarray = find_maximun_subarray(array, low, mid);
	auto right_maximun_subarray = find_maximun_subarray(array, mid, high);
	auto mid_maximun_subarray = find_max_crossing_subarray(array, low, mid, high);

	if (std::get<2>(left_maximun_subarray) >= std::get<2>(right_maximun_subarray) && 
		std::get<2>(left_maximun_subarray) >= std::get<2>(mid_maximun_subarray)) {
		return left_maximun_subarray;
	} else if (std::get<2>(right_maximun_subarray) >= std::get<2>(left_maximun_subarray) && 
			   std::get<2>(right_maximun_subarray) >= std::get<2>(mid_maximun_subarray)) {
		return right_maximun_subarray;
	} else {
		return mid_maximun_subarray;
	}
}

int main() {
	int array[] = { 3, 5, 2, -1, -8, 10 };
	auto maximun_subarray = find_maximun_subarray(array, 0, sizeof(array) / sizeof(int));
	std::cout << std::get<0>(maximun_subarray) << '\n';
	std::cout << std::get<1>(maximun_subarray) << '\n';
	std::cout << std::get<2>(maximun_subarray) << '\n';
	return 0;
}