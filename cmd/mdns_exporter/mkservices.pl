#!/usr/bin/perl
use strict;

# horrible hacky script to generate a map of service names
# from /etc/services

my $registry = `curl https://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.xml` or die;

my @lines = split /\n/, $registry;

my %services;
my ($name, $has_name, $descr, $has_descr) = ("", 0, "", 0);

for my $l (@lines) {
	# the first person, bail out
	if ($l =~ q|<person|) {
		last;
	}

	if ($l =~ q|<name>([^<]+)</name>|) {
		$name = $1;
		$has_name = 1;
	}
	
	# yeah this loses services with multiline descriptions and
	# I don't care
	
	if ($l =~ q|<description>([^<]+)</description>|) {
		$descr = $1;
		$has_descr = 1;
	}
	
	if ($has_name && $has_descr) {
		$services{$name} = $descr;
		
		$has_name = 0;
		$has_descr = 0;
	}
}

print qq|
package main

var serviceDescriptions map[string]string = map[string]string{
|;

my @keys = keys(%services);
@keys = sort(@keys);

for my $k (@keys) {
	my $value = $services{$k};
	$value =~ s/"/\\"/g;
	print qq|\t"_$k": "$value",\n|;
}

print qq|
}
|
