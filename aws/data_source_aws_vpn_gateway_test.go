package aws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAwsVpnGateway_unattached(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsVpnGatewayUnattachedConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.aws_vpn_gateway.test_by_id", "id",
						"aws_vpn_gateway.unattached", "id"),
					resource.TestCheckResourceAttrPair(
						"data.aws_vpn_gateway.test_by_tags", "id",
						"aws_vpn_gateway.unattached", "id"),
					resource.TestCheckResourceAttrPair(
						"data.aws_vpn_gateway.test_by_amazon_side_asn", "id",
						"aws_vpn_gateway.unattached", "id"),
					resource.TestCheckResourceAttrSet("data.aws_vpn_gateway.test_by_id", "state"),
					resource.TestCheckResourceAttr("data.aws_vpn_gateway.test_by_tags", "tags.%", "3"),
					resource.TestCheckNoResourceAttr("data.aws_vpn_gateway.test_by_id", "attached_vpc_id"),
					resource.TestCheckResourceAttr("data.aws_vpn_gateway.test_by_amazon_side_asn", "amazon_side_asn", "4294967293"),
				),
			},
		},
	})
}

func TestAccDataSourceAwsVpnGateway_attached(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsVpnGatewayAttachedConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.aws_vpn_gateway.test_by_attached_vpc_id", "id",
						"aws_vpn_gateway.attached", "id"),
					resource.TestCheckResourceAttrPair(
						"data.aws_vpn_gateway.test_by_attached_vpc_id", "attached_vpc_id",
						"aws_vpc.foo", "id"),
					resource.TestMatchResourceAttr("data.aws_vpn_gateway.test_by_attached_vpc_id", "state", regexp.MustCompile("(?i)available")),
				),
			},
		},
	})
}

func testAccDataSourceAwsVpnGatewayUnattachedConfig(rInt int) string {
	return fmt.Sprintf(`
resource "aws_vpn_gateway" "unattached" {
  tags = {
    Name = "terraform-testacc-vpn-gateway-data-source-unattached-%d"
    ABC  = "testacc-%d"
    XYZ  = "testacc-%d"
  }

  amazon_side_asn = 4294967293
}

data "aws_vpn_gateway" "test_by_id" {
  id = "${aws_vpn_gateway.unattached.id}"
}

data "aws_vpn_gateway" "test_by_tags" {
  tags = "${aws_vpn_gateway.unattached.tags}"
}

data "aws_vpn_gateway" "test_by_amazon_side_asn" {
  amazon_side_asn = "${aws_vpn_gateway.unattached.amazon_side_asn}"
  state           = "available"
}
`, rInt, rInt+1, rInt-1)
}

func testAccDataSourceAwsVpnGatewayAttachedConfig(rInt int) string {
	return fmt.Sprintf(`
resource "aws_vpc" "foo" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = "terraform-testacc-vpn-gateway-data-source-attached-%d"
  }
}

resource "aws_vpn_gateway" "attached" {
  tags = {
    Name = "terraform-testacc-vpn-gateway-data-source-attached-%d"
  }
}

resource "aws_vpn_gateway_attachment" "vpn_attachment" {
  vpc_id         = "${aws_vpc.foo.id}"
  vpn_gateway_id = "${aws_vpn_gateway.attached.id}"
}

data "aws_vpn_gateway" "test_by_attached_vpc_id" {
  attached_vpc_id = "${aws_vpn_gateway_attachment.vpn_attachment.vpc_id}"
}
`, rInt, rInt)
}
