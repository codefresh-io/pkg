// Copyright 2021 The Codefresh Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.// Code generated by interfacer; DO NOT EDIT

package github

import (
	"context"
	"github.com/google/go-github/v32/github"
	"io"
	"net/http"
	"net/url"
	"os"
)

// Repositories is an interface generated for "github.com/google/go-github/v32/github.RepositoriesService".
type Repositories interface {
	AddAdminEnforcement(context.Context, string, string, string) (*github.AdminEnforcement, *github.Response, error)
	AddAppRestrictions(context.Context, string, string, string, []string) ([]*github.App, *github.Response, error)
	AddCollaborator(context.Context, string, string, string, *github.RepositoryAddCollaboratorOptions) (*github.CollaboratorInvitation, *github.Response, error)
	CompareCommits(context.Context, string, string, string, string) (*github.CommitsComparison, *github.Response, error)
	Create(context.Context, string, *github.Repository) (*github.Repository, *github.Response, error)
	CreateComment(context.Context, string, string, string, *github.RepositoryComment) (*github.RepositoryComment, *github.Response, error)
	CreateDeployment(context.Context, string, string, *github.DeploymentRequest) (*github.Deployment, *github.Response, error)
	CreateDeploymentStatus(context.Context, string, string, int64, *github.DeploymentStatusRequest) (*github.DeploymentStatus, *github.Response, error)
	CreateFile(context.Context, string, string, string, *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error)
	CreateFork(context.Context, string, string, *github.RepositoryCreateForkOptions) (*github.Repository, *github.Response, error)
	CreateFromTemplate(context.Context, string, string, *github.TemplateRepoRequest) (*github.Repository, *github.Response, error)
	CreateHook(context.Context, string, string, *github.Hook) (*github.Hook, *github.Response, error)
	CreateKey(context.Context, string, string, *github.Key) (*github.Key, *github.Response, error)
	CreateProject(context.Context, string, string, *github.ProjectOptions) (*github.Project, *github.Response, error)
	CreateRelease(context.Context, string, string, *github.RepositoryRelease) (*github.RepositoryRelease, *github.Response, error)
	CreateStatus(context.Context, string, string, string, *github.RepoStatus) (*github.RepoStatus, *github.Response, error)
	Delete(context.Context, string, string) (*github.Response, error)
	DeleteComment(context.Context, string, string, int64) (*github.Response, error)
	DeleteDeployment(context.Context, string, string, int64) (*github.Response, error)
	DeleteFile(context.Context, string, string, string, *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error)
	DeleteHook(context.Context, string, string, int64) (*github.Response, error)
	DeleteInvitation(context.Context, string, string, int64) (*github.Response, error)
	DeleteKey(context.Context, string, string, int64) (*github.Response, error)
	DeletePreReceiveHook(context.Context, string, string, int64) (*github.Response, error)
	DeleteRelease(context.Context, string, string, int64) (*github.Response, error)
	DeleteReleaseAsset(context.Context, string, string, int64) (*github.Response, error)
	DisableAutomatedSecurityFixes(context.Context, string, string) (*github.Response, error)
	DisableDismissalRestrictions(context.Context, string, string, string) (*github.PullRequestReviewsEnforcement, *github.Response, error)
	DisablePages(context.Context, string, string) (*github.Response, error)
	DisableVulnerabilityAlerts(context.Context, string, string) (*github.Response, error)
	Dispatch(context.Context, string, string, github.DispatchRequestOptions) (*github.Repository, *github.Response, error)
	DownloadContents(context.Context, string, string, string, *github.RepositoryContentGetOptions) (io.ReadCloser, error)
	DownloadReleaseAsset(context.Context, string, string, int64, *http.Client) (io.ReadCloser, string, error)
	Edit(context.Context, string, string, *github.Repository) (*github.Repository, *github.Response, error)
	EditHook(context.Context, string, string, int64, *github.Hook) (*github.Hook, *github.Response, error)
	EditRelease(context.Context, string, string, int64, *github.RepositoryRelease) (*github.RepositoryRelease, *github.Response, error)
	EditReleaseAsset(context.Context, string, string, int64, *github.ReleaseAsset) (*github.ReleaseAsset, *github.Response, error)
	EnableAutomatedSecurityFixes(context.Context, string, string) (*github.Response, error)
	EnablePages(context.Context, string, string, *github.Pages) (*github.Pages, *github.Response, error)
	EnableVulnerabilityAlerts(context.Context, string, string) (*github.Response, error)
	Get(context.Context, string, string) (*github.Repository, *github.Response, error)
	GetAdminEnforcement(context.Context, string, string, string) (*github.AdminEnforcement, *github.Response, error)
	GetArchiveLink(context.Context, string, string, github.ArchiveFormat, *github.RepositoryContentGetOptions, bool) (*url.URL, *github.Response, error)
	GetBranch(context.Context, string, string, string) (*github.Branch, *github.Response, error)
	GetBranchProtection(context.Context, string, string, string) (*github.Protection, *github.Response, error)
	GetByID(context.Context, int64) (*github.Repository, *github.Response, error)
	GetCodeOfConduct(context.Context, string, string) (*github.CodeOfConduct, *github.Response, error)
	GetCombinedStatus(context.Context, string, string, string, *github.ListOptions) (*github.CombinedStatus, *github.Response, error)
	GetComment(context.Context, string, string, int64) (*github.RepositoryComment, *github.Response, error)
	GetCommit(context.Context, string, string, string) (*github.RepositoryCommit, *github.Response, error)
	GetCommitRaw(context.Context, string, string, string, github.RawOptions) (string, *github.Response, error)
	GetCommitSHA1(context.Context, string, string, string, string) (string, *github.Response, error)
	GetCommunityHealthMetrics(context.Context, string, string) (*github.CommunityHealthMetrics, *github.Response, error)
	GetContents(context.Context, string, string, string, *github.RepositoryContentGetOptions) (*github.RepositoryContent, []*github.RepositoryContent, *github.Response, error)
	GetDeployment(context.Context, string, string, int64) (*github.Deployment, *github.Response, error)
	GetDeploymentStatus(context.Context, string, string, int64, int64) (*github.DeploymentStatus, *github.Response, error)
	GetHook(context.Context, string, string, int64) (*github.Hook, *github.Response, error)
	GetKey(context.Context, string, string, int64) (*github.Key, *github.Response, error)
	GetLatestPagesBuild(context.Context, string, string) (*github.PagesBuild, *github.Response, error)
	GetLatestRelease(context.Context, string, string) (*github.RepositoryRelease, *github.Response, error)
	GetPageBuild(context.Context, string, string, int64) (*github.PagesBuild, *github.Response, error)
	GetPagesInfo(context.Context, string, string) (*github.Pages, *github.Response, error)
	GetPermissionLevel(context.Context, string, string, string) (*github.RepositoryPermissionLevel, *github.Response, error)
	GetPreReceiveHook(context.Context, string, string, int64) (*github.PreReceiveHook, *github.Response, error)
	GetPullRequestReviewEnforcement(context.Context, string, string, string) (*github.PullRequestReviewsEnforcement, *github.Response, error)
	GetReadme(context.Context, string, string, *github.RepositoryContentGetOptions) (*github.RepositoryContent, *github.Response, error)
	GetRelease(context.Context, string, string, int64) (*github.RepositoryRelease, *github.Response, error)
	GetReleaseAsset(context.Context, string, string, int64) (*github.ReleaseAsset, *github.Response, error)
	GetReleaseByTag(context.Context, string, string, string) (*github.RepositoryRelease, *github.Response, error)
	GetRequiredStatusChecks(context.Context, string, string, string) (*github.RequiredStatusChecks, *github.Response, error)
	GetSignaturesProtectedBranch(context.Context, string, string, string) (*github.SignaturesProtectedBranch, *github.Response, error)
	GetVulnerabilityAlerts(context.Context, string, string) (bool, *github.Response, error)
	IsCollaborator(context.Context, string, string, string) (bool, *github.Response, error)
	License(context.Context, string, string) (*github.RepositoryLicense, *github.Response, error)
	List(context.Context, string, *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error)
	ListAll(context.Context, *github.RepositoryListAllOptions) ([]*github.Repository, *github.Response, error)
	ListAllTopics(context.Context, string, string) ([]string, *github.Response, error)
	ListApps(context.Context, string, string, string) ([]*github.App, *github.Response, error)
	ListBranches(context.Context, string, string, *github.BranchListOptions) ([]*github.Branch, *github.Response, error)
	ListBranchesHeadCommit(context.Context, string, string, string) ([]*github.BranchCommit, *github.Response, error)
	ListByOrg(context.Context, string, *github.RepositoryListByOrgOptions) ([]*github.Repository, *github.Response, error)
	ListCodeFrequency(context.Context, string, string) ([]*github.WeeklyStats, *github.Response, error)
	ListCollaborators(context.Context, string, string, *github.ListCollaboratorsOptions) ([]*github.User, *github.Response, error)
	ListComments(context.Context, string, string, *github.ListOptions) ([]*github.RepositoryComment, *github.Response, error)
	ListCommitActivity(context.Context, string, string) ([]*github.WeeklyCommitActivity, *github.Response, error)
	ListCommitComments(context.Context, string, string, string, *github.ListOptions) ([]*github.RepositoryComment, *github.Response, error)
	ListCommits(context.Context, string, string, *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error)
	ListContributors(context.Context, string, string, *github.ListContributorsOptions) ([]*github.Contributor, *github.Response, error)
	ListContributorsStats(context.Context, string, string) ([]*github.ContributorStats, *github.Response, error)
	ListDeploymentStatuses(context.Context, string, string, int64, *github.ListOptions) ([]*github.DeploymentStatus, *github.Response, error)
	ListDeployments(context.Context, string, string, *github.DeploymentsListOptions) ([]*github.Deployment, *github.Response, error)
	ListForks(context.Context, string, string, *github.RepositoryListForksOptions) ([]*github.Repository, *github.Response, error)
	ListHooks(context.Context, string, string, *github.ListOptions) ([]*github.Hook, *github.Response, error)
	ListInvitations(context.Context, string, string, *github.ListOptions) ([]*github.RepositoryInvitation, *github.Response, error)
	ListKeys(context.Context, string, string, *github.ListOptions) ([]*github.Key, *github.Response, error)
	ListLanguages(context.Context, string, string) (map[string]int, *github.Response, error)
	ListPagesBuilds(context.Context, string, string, *github.ListOptions) ([]*github.PagesBuild, *github.Response, error)
	ListParticipation(context.Context, string, string) (*github.RepositoryParticipation, *github.Response, error)
	ListPreReceiveHooks(context.Context, string, string, *github.ListOptions) ([]*github.PreReceiveHook, *github.Response, error)
	ListProjects(context.Context, string, string, *github.ProjectListOptions) ([]*github.Project, *github.Response, error)
	ListPunchCard(context.Context, string, string) ([]*github.PunchCard, *github.Response, error)
	ListReleaseAssets(context.Context, string, string, int64, *github.ListOptions) ([]*github.ReleaseAsset, *github.Response, error)
	ListReleases(context.Context, string, string, *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error)
	ListRequiredStatusChecksContexts(context.Context, string, string, string) ([]string, *github.Response, error)
	ListStatuses(context.Context, string, string, string, *github.ListOptions) ([]*github.RepoStatus, *github.Response, error)
	ListTags(context.Context, string, string, *github.ListOptions) ([]*github.RepositoryTag, *github.Response, error)
	ListTeams(context.Context, string, string, *github.ListOptions) ([]*github.Team, *github.Response, error)
	ListTrafficClones(context.Context, string, string, *github.TrafficBreakdownOptions) (*github.TrafficClones, *github.Response, error)
	ListTrafficPaths(context.Context, string, string) ([]*github.TrafficPath, *github.Response, error)
	ListTrafficReferrers(context.Context, string, string) ([]*github.TrafficReferrer, *github.Response, error)
	ListTrafficViews(context.Context, string, string, *github.TrafficBreakdownOptions) (*github.TrafficViews, *github.Response, error)
	Merge(context.Context, string, string, *github.RepositoryMergeRequest) (*github.RepositoryCommit, *github.Response, error)
	OptionalSignaturesOnProtectedBranch(context.Context, string, string, string) (*github.Response, error)
	PingHook(context.Context, string, string, int64) (*github.Response, error)
	RemoveAdminEnforcement(context.Context, string, string, string) (*github.Response, error)
	RemoveAppRestrictions(context.Context, string, string, string, []string) ([]*github.App, *github.Response, error)
	RemoveBranchProtection(context.Context, string, string, string) (*github.Response, error)
	RemoveCollaborator(context.Context, string, string, string) (*github.Response, error)
	RemovePullRequestReviewEnforcement(context.Context, string, string, string) (*github.Response, error)
	ReplaceAllTopics(context.Context, string, string, []string) ([]string, *github.Response, error)
	ReplaceAppRestrictions(context.Context, string, string, string, []string) ([]*github.App, *github.Response, error)
	RequestPageBuild(context.Context, string, string) (*github.PagesBuild, *github.Response, error)
	RequireSignaturesOnProtectedBranch(context.Context, string, string, string) (*github.SignaturesProtectedBranch, *github.Response, error)
	TestHook(context.Context, string, string, int64) (*github.Response, error)
	Transfer(context.Context, string, string, github.TransferRequest) (*github.Repository, *github.Response, error)
	UpdateBranchProtection(context.Context, string, string, string, *github.ProtectionRequest) (*github.Protection, *github.Response, error)
	UpdateComment(context.Context, string, string, int64, *github.RepositoryComment) (*github.RepositoryComment, *github.Response, error)
	UpdateFile(context.Context, string, string, string, *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error)
	UpdateInvitation(context.Context, string, string, int64, string) (*github.RepositoryInvitation, *github.Response, error)
	UpdatePages(context.Context, string, string, *github.PagesUpdate) (*github.Response, error)
	UpdatePreReceiveHook(context.Context, string, string, int64, *github.PreReceiveHook) (*github.PreReceiveHook, *github.Response, error)
	UpdatePullRequestReviewEnforcement(context.Context, string, string, string, *github.PullRequestReviewsEnforcementUpdate) (*github.PullRequestReviewsEnforcement, *github.Response, error)
	UpdateRequiredStatusChecks(context.Context, string, string, string, *github.RequiredStatusChecksRequest) (*github.RequiredStatusChecks, *github.Response, error)
	UploadReleaseAsset(context.Context, string, string, int64, *github.UploadOptions, *os.File) (*github.ReleaseAsset, *github.Response, error)
}
