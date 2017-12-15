#include <iostream>

int main(int argc,char* argv[])
{
    if (argc != 3)
    {
        std::cerr << "argc is not correct!" << std::endl;
        return 0;
    }

    std::cout << argv[1] << std::endl << argv[2] << std::endl;
    return 0;
}
