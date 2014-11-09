require 'pathname'
require 'mkmf'
require 'colorize'

# some standard folders
$test_folder   = 'tests'
$mock_folder   = "#{$test_folder}/mocks"
$vendor_folder = 'vendor'
$vendor_bin    = "#{$vendor_folder}/bin"
$vendor_src    = "#{$vendor_folder}/src"
$go_common     = "#{$vendor_src}/go.classmarkets.com/go-common"

# general settings
$mock_package_name = 'mocks'

$vendor_bins = {
    'ginkgo'  => 'github.com/onsi/ginkgo/ginkgo',
    'mockgen' => 'code.google.com/p/gomock/mockgen'
}

$mocks = {
    'shell.go' => {},
    'task.go' => {},
    'file_system.go' => {
        imports: '.=github.com/fgrosse/grobot'
    }
}

task :default => :test

directory $test_folder
directory $mock_folder
directory $vendor_folder

desc 'Run all ginkgo tests'
task :test => [$test_folder, :mocks, "#{$vendor_bin}/ginkgo"] do
    #sh "ginkgo --noisyPendings=false -r #{$test_folder}"
end

task :mocks => [$mock_folder, "#{$vendor_bin}/mockgen"] + $mocks.keys do
    $mocks.each do |mock_source, conf|
        mock_file = mock_file_name(mock_source, conf)
        puts "Generating mock #{mock_file} from #{mock_source}".yellow
        command = "#{$vendor_bin}/mockgen -source #{mock_source} -destination #{mock_file} -package #{$mock_package_name}"
        if conf.has_key? :imports
            command += " -imports \"#{conf[:imports]}\""
        end
        sh command
    end
end

# Generic rule to build vendor binaries
rule /^#{$vendor_folder}\/bin\/\w+/ => ->(f) {vendor_source_for(f)} do |file|
    vendor_package = /^#{$vendor_src}\/(.+)/.match(file.source)[1]
    puts "Compiling #{file.name}..".yellow
    sh "go build -o #{file.name} #{vendor_package}"
end

# Get the source folder for a vendor binary from the global $vendor_binaries hash
def vendor_source_for(f)
    bin_name = /^#{$vendor_bin}\/(\w+)/.match(f)[1]
    source_missing_hint = <<eos

No source for vendor bin '#{bin_name}' defined in $vendor_bin
Please provide the source to compile #{$vendor_bin}/#{bin_name} in the format:

    $vendor_bins = {
      'ginkgo'  => 'github.com/onsi/ginkgo/ginkgo'
    }

eos

    raise source_missing_hint unless $vendor_bins.has_key? bin_name
    return "#{$vendor_src}/#{$vendor_bins[bin_name]}"
end

# Generic rule to build mocks
rule /^#{$mock_folder}\/\w\.go$/ => ->(f) {mock_source_for(f)} do
    
end

# Get the source of a mock from the global $mocks hash
def mock_source_for(mock_file) 
    base_name = Pathname.new(mock_file).basename
    base_name.gsub!(/_mock\.go$/, '.go')
    puts "#{mock_file} => #{base_name}"
end

def mock_file_name(mock_source, conf)
    if conf.has_key? :mock_file_name
       mock_file = conf[:mock_file_name]
    else
       base_name = Pathname.new(mock_source).basename.to_s
       mock_file = base_name.gsub(/\.go$/, '') + '_mock.go'
    end
    return "#{$mock_folder}/#{mock_file}"
end

namespace :debug do
    desc 'Print out all defined mocks for this project'
    task :mocks do
        $mocks.each do |mock_source, conf|
            mock_file = mock_file_name(mock_source, conf) 
            puts "#{mock_file}: ./#{mock_source}"
            puts "    config: #{conf.inspect}\n"
        end
    end

    desc 'Print out all used folders'
    task :folders do
        puts "./#{$test_folder}"
        puts "./#{$mock_folder}"
        puts "./#{$vendor_folder}"
        puts "./#{$vendor_src}"
        puts "./#{$vendor_bin}"
    end
end
