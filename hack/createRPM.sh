#!/bin/bash
# Based on https://github.com/cernbox/ocmd/blob/master/Makefile

set -xe

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SPECFILE="${SCRIPTDIR}/inlets.spec"
PACKAGE=$(awk '$1 == "Name:"     { print $2 }' ${SPECFILE})
VERSION=$(awk '$1 == "Version:"  { print $2 }' ${SPECFILE})
RELEASE=$(awk '$1 == "Release:"  { print $2 }' ${SPECFILE})

rpmbuild=$(mktemp -d)

mkdir -p ${rpmbuild}/{RPMS,SPECS,SOURCES,BUILD,SRPMS}
mkdir -p ${rpmbuild}/RPMS/x86_64/

dep ensure
go build

TEMPDIR=$(mktemp -d)
mkdir -p ${TEMPDIR}/${PACKAGE}-${VERSION}

cp -r hack/inlets.service inlets LICENSE README.md ${TEMPDIR}/${PACKAGE}-${VERSION}/
tar cpfz ${rpmbuild}/SOURCES/${PACKAGE}-${VERSION}.tar.gz --directory=${TEMPDIR} ${PACKAGE}-${VERSION}

cp ${SPECFILE} ${rpmbuild}/SPECS/

rpmbuild --define="_topdir ${rpmbuild}" \
        --define="_sourcedir %{_topdir}/SOURCES" \
        --define="_builddir %{_topdir}/BUILD" \
        --define="_srcrpmdir %{_topdir}/SRPMS" \
        --define="_rpmdir %{_topdir}/RPMS" \
        --nodeps -bb ${rpmbuild}/SPECS/${PACKAGE}.spec

cp ${rpmbuild}/RPMS/x86_64/* .

rm -Rf ${rpmbuild} ${TEMPDIR}
