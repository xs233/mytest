#include <iostream>

void recursive_activity_selector(int start[], int stop[], int s_loc, int size) {
	int index = s_loc + 1;
	while (index < size) {
		if (stop[s_loc] <= start[index]) {
			break;
		}
		++index;
	}

	if (index < size) {
		std::cout << index <<'\n';
		recursive_activity_selector(start, stop, index, size);
	}
}

void greedy_activity_selector(int start[], int stop[], int size) {
	int s_loc = 0;
	for (int i = 1; i < size; ++i) {
		if (stop[s_loc] <= start[i]) {
			s_loc = i;
			std::cout << s_loc << '\n';
		}
	}
}

int main() {
	int start[] = { 1, 3, 0, 5, 3, 5, 6, 8, 8, 2, 12 };
	int stop[] = { 4, 5, 6, 7, 9, 9, 10, 11, 12, 14, 16 };
	recursive_activity_selector(start, stop, 0, sizeof(start) / sizeof(int));
	greedy_activity_selector(start, stop, sizeof(start) / sizeof(int));
	return 0;
}