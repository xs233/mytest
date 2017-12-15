#include <boost/archive/iterators/base64_from_binary.hpp>
#include <boost/archive/iterators/binary_from_base64.hpp>
#include <boost/archive/iterators/insert_linebreaks.hpp>
#include <boost/archive/iterators/transform_width.hpp>
#include <boost/archive/iterators/ostream_iterator.hpp>
#include <sstream>
#include <string>
#include <iostream>

int main()
{
    using namespace boost::archive::iterators;
    // encode
    std::string test = "xushen";

    std::stringstream os;
    typedef 
        insert_linebreaks<         // insert line breaks every 72 characters
            base64_from_binary<    // convert binary values to base64 characters
                transform_width<   // retrieve 6 bit integers from a sequence of 8 bit bytes
                    const char *,
                    6,
                    8
                >
            > 
            ,72
        > 
        base64_text; // compose all the above operations in to a new iterator

    std::copy(
        base64_text(test.c_str()),
        base64_text(test.c_str() + test.size()),
        ostream_iterator<char>(os)
    );

    std::cout << os.str() << '\n';
    std::cout << os.str().size() << '\n';

    // decode
    const char* de_str = "eHVzaGVu";
    typedef 
        insert_linebreaks<    
            transform_width<     
                binary_from_base64<       
                    const char *
                    >,
                8,
                6
            >
            ,72
        > 
        debase64_text; // compose all the above operations in to a new iterator

    std::copy(
        debase64_text(de_str),
        debase64_text(de_str + 8),
        ostream_iterator<char>(os)
    );
    std::cout << os.str() << '\n';
    std::cout << os.str().size() << '\n';

    // for (auto it = debase64_text(de_str.c_str()); it != debase64_text(de_str.c_str() + 8); ++it) {
    //     std::cout << *it << '\n';
    // }

    // std::string decode(debase64_text(de_str.c_str()), debase64_text(de_str.c_str() + 8));
    // std::cout << decode << '\n';
    return 0;
}