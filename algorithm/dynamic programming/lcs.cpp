#include <cstring>
#include <iostream>

enum DIRECTION {
	UP,
	LEFT,
	DECUMBENT
};

struct status {
	int max_len;
	DIRECTION dir;
};

int get_len(status** status_array, int i, int j) {
	if (i < 0) {
		i = 0;
	}

	if (j < 0) {
		j = 0;
	}

	return status_array[i][j].max_len;
}

status** lcs(char x[], int x_len, char y[], int y_len) {
	status** status_array = new status*[x_len];
	for (int i = 0; i < x_len; ++i) {
		status_array[i] = new status[y_len];
		std::memset(status_array[i], '\0', y_len * sizeof(status));
		for (int j = 0; j < y_len; ++j) {
			if (x[i] == y[j]) {
				status_array[i][j].max_len = get_len(status_array, i - 1, j - 1) + 1;
				status_array[i][j].dir = DECUMBENT;
			} else if (get_len(status_array, i, j - 1) >= get_len(status_array, i - 1, j)) {
				status_array[i][j].max_len = get_len(status_array, i, j - 1);
				status_array[i][j].dir = UP;
			} else {
				status_array[i][j].max_len = get_len(status_array, i - 1, j);
				status_array[i][j].dir = LEFT;
			}
		}
	}
	return status_array;
}

void print_lcs(status** status_array, char x[], int x_len, int y_len) {
	if (x_len < 0 && y_len < 0) {
		return;
	}

	switch (status_array[x_len][y_len].dir) {
	case UP:
		print_lcs(status_array, x, x_len, y_len - 1);
		break;
	case LEFT:
		print_lcs(status_array, x, x_len - 1, y_len);
		break;
	case DECUMBENT:
		print_lcs(status_array, x, x_len - 1, y_len - 1);
		std::cout << x[x_len] << '\n';
		break;
	}
}

void wrap_print_lcs(status** status_array, char x[], int x_len, int y_len) {
	print_lcs(status_array, x, --x_len, --y_len);
}

int main() {
	char x[] = { 'A', 'B', 'C', 'B', 'D', 'A', 'B' };
	char y[] = { 'B', 'D', 'C', 'A', 'B', 'A' };
	status** status_array = lcs(x, sizeof(x) / sizeof(char), y, sizeof(y) / sizeof(char));
	wrap_print_lcs(status_array, x, sizeof(x) / sizeof(char), sizeof(y) / sizeof(char));
	for (int i = 0; i < sizeof(x) / sizeof(char); ++i) {
		delete[] status_array[i];
	}
	delete[] status_array;
	return 0;
}