%if 0%{?fedora}
%global with_bundled 0
%global with_debug 1
%global with_check 1
%global with_unit_test 0
%else
%global with_bundled 1
%global with_debug 1
# no test files so far
%global with_check 0
# no test files so far
%global with_unit_test 0
%endif

%if 0%{?with_debug}
%global _dwz_low_mem_die_limit 0
%else
%global debug_package   %{nil}
%endif

%global provider        github
%global provider_tld    com
%global project         stefwalter
%global repo            oci-kvm-hook
# https://github.com/stefwalter/oci-kvm-hook
%global provider_prefix %{provider}.%{provider_tld}/%{project}/%{repo}
%global import_path     %{provider_prefix}

Name:           %{repo}
Version:        0.2
Release:        2%{?dist}
Summary:        Golang binary to mount /dev/kvm into OCI containers
License:        ASL 2.0
URL:            https://%{import_path}
Source0:        https://%{import_path}/archive/%{version}.tar.gz

# If go_compiler is not set to 1, there is no virtual provide. Use golang instead.
BuildRequires:  %{?go_compiler:compiler(go-compiler)}%{!?go_compiler:golang}

%if ! 0%{?with_bundled}
BuildRequires: golang(gopkg.in/yaml.v1)
%endif
BuildRequires:   go-md2man

ExclusiveArch:  %{?go_arches:%{go_arches}}%{!?go_arches:%{ix86} x86_64 aarch64 %{arm}}

%description
%{summary}

%prep
%setup -q -n %{repo}-%{version}

%build
%if ! 0%{?with_bundled}
export GOPATH=$(pwd):%{gopath}
%else
export GOPATH=$(pwd):$(pwd)/Godeps/_workspace:%{gopath}
%endif

%make_build

%install
%make_install

%files
%license LICENSE
%doc %{name}.1.md README.md
%dir %{_libexecdir}/oci
%dir %{_libexecdir}/oci/hooks.d
%{_libexecdir}/oci/hooks.d/%{name}
%{_mandir}/man1/%{name}.1*

%changelog
* Thu Sep 21 2017 Stef Walter <stefw@redhat.com> - 0.2-2
- Updated for package review

* Wed Sep 20 2017 Stef Walter <stefw@redhat.com> - 0.2-1
- Copy /dev/kvm permissions from host
- Avoid nsenter --cgroup option for compatibility

* Wed Sep 20 2017 Stef Walter <stefw@redhat.com> - 0.1-1
- Initial release
