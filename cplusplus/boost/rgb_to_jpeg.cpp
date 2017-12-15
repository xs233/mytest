#include <boost/gil/extension/io/jpeg_io.hpp>

const unsigned width  = 320;
const unsigned height = 200;

// Raw data.
unsigned char r[width * height];  // red
unsigned char g[width * height];  // green
unsigned char b[width * height];  // blue

int main()
{
    boost::gil::rgb8c_planar_view_t view = boost::gil::planar_rgb_view(width, height, r, g, b, width);
    boost::gil::jpeg_write_view("out.jpg", view);
    return 0;
}