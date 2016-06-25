#include <cerrno>
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <iostream>

#include <boost/program_options.hpp>
#include <boost/lexical_cast.hpp>

struct CompilerOptions
{
  bool verbose;
  bool help;
  CompilerOptions(): verbose(false), help(false) { }
};

CompilerOptions options;

int main(int args, char *argv[])
{
    boost::program_options::options_description desc("Options");
    desc.add_options()
        ("help", "Options related to the program.")
        ("verbose,v", boost::program_options::bool_switch(&options.verbose)->default_value(false), "Print to stdout information as job is processed.")
        ;

  // parse command line options
  boost::program_options::variables_map vm;
  try
  {
      boost::program_options::store(boost::program_options::parse_command_line(args, argv, desc), vm);
      boost::program_options::notify(vm);
  }
  catch(std::exception &e)
  { 
    std::cout << e.what() << std::endl;
    return EXIT_FAILURE;
  }

  // check help flag
  if (vm.count("help"))
  {
    std::cout << desc << std::endl;
    return EXIT_SUCCESS;
  }

  if (options.verbose) {
    // std::cerr << "Initializing worker..." << std::endl;
    // return EXIT_FAILURE;
  }
  return EXIT_SUCCESS;
}

