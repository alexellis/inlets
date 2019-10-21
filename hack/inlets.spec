# Based on https://github.com/cernbox/ocmd/blob/master/ocmd.spec
%global debug_package %{nil}

Name:           inlets
Release:        1%{?dist}
Summary:        Expose your local endpoints to the Internet
Version:        2.6.1
License:        MIT
URL:            https://github.com/inlets/inlets
Source0:        https://github.com/inlets/inlets/archive/%{name}-%{version}.tar.gz

BuildRequires: systemd
BuildRequires: systemd-rpm-macros

%description
Expose your local endpoints to the Internet

%prep
%setup -n %{name}-%{version}

%build

%install
install -d %{buildroot}/usr/local/bin
install -d %{buildroot}/etc/default
install -d %{buildroot}%{_unitdir}
install -p -m 0755 inlets %{buildroot}/usr/local/bin/inlets
install -p -m 644 inlets.service %{buildroot}%{_unitdir}

%files
%license LICENSE
%doc README.md
/usr/local/bin/inlets
%{_unitdir}/*

%post
%systemd_post inlets.service

%preun
%systemd_preun inlets.service

%postun
%systemd_postun inlets.service

%changelog
* Mon Oct 21 2019 Eduardo Minguez Perez <e.minguez@gmail.com> - 2.6.1-1
- Initial package
