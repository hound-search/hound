require 'fileutils'

GOPATH=File.open('.gaan').map do |x|
  File.absolute_path(File.join(File.dirname(__FILE__), x.strip))
end.join(':')

ENV.update({
  'GOPATH' => GOPATH
})

file 'bin/hound' => FileList['src/hound/**/*', 'src/ansi/**/*'] do
  host = ENV['HOST'] || 'localhost:6080'
  args = ['go', 'build', '-o', 'bin/hound', '-ldflags', "-X main.defaultHost #{host}"] + FileList['src/hound/cmds/hound/*.go']
  sh *args
end

file 'bin/houndd' => FileList['src/hound/**/*'] do
  sh 'go', 'build', '-o', 'bin/houndd', 'src/hound/cmds/houndd/main.go'
end

task :default => [
  'bin/houndd',
  'bin/hound',
]

task :clean do
  FileUtils::rm_rf('bin')
end

task :test do
  pkgs = FileList['src/**/*_test.go'].map do |x|
    File.dirname(x[4,x.length])
  end.uniq
  args = ['go', 'test'].concat(pkgs)
  sh *args
end
