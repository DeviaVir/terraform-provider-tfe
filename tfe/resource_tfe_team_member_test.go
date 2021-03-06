package tfe

import (
	"fmt"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestPackTeamMemberID(t *testing.T) {
	cases := []struct {
		team string
		user string
		id   string
	}{
		{
			team: "team-47qC3LmA47piVan7",
			user: "sander",
			id:   "team-47qC3LmA47piVan7/sander",
		},
	}

	for _, tc := range cases {
		id := packTeamMemberID(tc.team, tc.user)

		if tc.id != id {
			t.Fatalf("expected ID %q, got %q", tc.id, id)
		}
	}

}

func TestUnpackTeamMemberID(t *testing.T) {
	cases := []struct {
		id   string
		team string
		user string
		err  bool
	}{
		{
			id:   "team-47qC3LmA47piVan7/sander",
			team: "team-47qC3LmA47piVan7",
			user: "sander",
			err:  false,
		},
		{
			id:   "team-47qC3LmA47piVan7|sander",
			team: "team-47qC3LmA47piVan7",
			user: "sander",
			err:  false,
		},
		{
			id:   "some-invalid-id",
			team: "",
			user: "",
			err:  true,
		},
	}

	for _, tc := range cases {
		team, user, err := unpackTeamMemberID(tc.id)
		if (err != nil) != tc.err {
			t.Fatalf("expected error is %t, got %v", tc.err, err)
		}

		if tc.team != team {
			t.Fatalf("expected team %q, got %q", tc.team, team)
		}

		if tc.user != user {
			t.Fatalf("expected user %q, got %q", tc.user, user)
		}
	}

}

func TestAccTFETeamMember_basic(t *testing.T) {
	user := &tfe.User{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTFETeamMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTFETeamMember_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTFETeamMemberExists(
						"tfe_team_member.foobar", user),
					testAccCheckTFETeamMemberAttributes(user),
					resource.TestCheckResourceAttr(
						"tfe_team_member.foobar", "username", "admin"),
				),
			},
		},
	})
}

func TestAccTFETeamMember_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTFETeamMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTFETeamMember_basic,
			},

			{
				ResourceName:      "tfe_team_member.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTFETeamMemberExists(
	n string, user *tfe.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		tfeClient := testAccProvider.Meta().(*tfe.Client)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		// Get the team ID and username.
		teamID, username, err := unpackTeamMemberID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error unpacking team member ID: %v", err)
		}

		users, err := tfeClient.TeamMembers.List(ctx, teamID)
		if err != nil && err != tfe.ErrResourceNotFound {
			return err
		}

		found := false
		for _, u := range users {
			if u.Username == username {
				found = true
				*user = *u
				break
			}
		}

		if !found {
			return fmt.Errorf("User not found")
		}

		return nil
	}
}

func testAccCheckTFETeamMemberAttributes(
	user *tfe.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if user.Username != "admin" {
			return fmt.Errorf("Bad username: %s", user.Username)
		}
		return nil
	}
}

func testAccCheckTFETeamMemberDestroy(s *terraform.State) error {
	tfeClient := testAccProvider.Meta().(*tfe.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tfe_team_member" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		// Get the team ID and username.
		teamID, username, err := unpackTeamMemberID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error unpacking team member ID: %v", err)
		}

		users, err := tfeClient.TeamMembers.List(ctx, teamID)
		if err != nil && err != tfe.ErrResourceNotFound {
			return err
		}

		found := false
		for _, u := range users {
			if u.Username == username {
				found = true
				break
			}
		}

		if found {
			return fmt.Errorf("User %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccTFETeamMember_basic = `
resource "tfe_organization" "foobar" {
  name  = "tst-terraform"
  email = "admin@company.com"
}

resource "tfe_team" "foobar" {
  name         = "team-test"
  organization = "${tfe_organization.foobar.id}"
}

resource "tfe_team_member" "foobar" {
  team_id  = "${tfe_team.foobar.id}"
  username = "admin"
}`
