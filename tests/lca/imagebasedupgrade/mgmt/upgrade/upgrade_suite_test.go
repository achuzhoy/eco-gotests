package upgrade_test

import (
	"testing"

	"runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	"github.com/openshift-kni/eco-gotests/tests/internal/reporter"
	"github.com/openshift-kni/eco-gotests/tests/lca/imagebasedupgrade/internal/seedimage"
	. "github.com/openshift-kni/eco-gotests/tests/lca/imagebasedupgrade/mgmt/internal/mgmtinittools"
	"github.com/openshift-kni/eco-gotests/tests/lca/imagebasedupgrade/mgmt/upgrade/internal/tsparams"
	_ "github.com/openshift-kni/eco-gotests/tests/lca/imagebasedupgrade/mgmt/upgrade/tests"
)

var _, currentFile, _, _ = runtime.Caller(0)

func TestUpgrade(t *testing.T) {
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = MGMTConfig.GetJunitReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Upgrade Suite", Label(tsparams.Labels...), reporterConfig)
}

var _ = BeforeSuite(func() {
	if MGMTConfig.SeedImageVersion == "" {
		seedInfo, err := seedimage.GetContent(APIClient, MGMTConfig.SeedImage)
		Expect(err).NotTo(HaveOccurred(), "error getting seed image info")

		MGMTConfig.SeedImageVersion = seedInfo.SeedClusterOCPVersion
	}
})

var _ = ReportAfterSuite("", func(report Report) {
	reportxml.Create(
		report, MGMTConfig.GetReportPath(), MGMTConfig.TCPrefix)
})

var _ = JustAfterEach(func() {
	reporter.ReportIfFailed(
		CurrentSpecReport(),
		currentFile,
		tsparams.ReporterNamespacesToDump,
		tsparams.ReporterCRDsToDump,
		clients.SetScheme)
})
