import http from 'k6/http';
import { check } from 'k6';

export const options = {
    vus: 1,
    duration: '60s',
    rps: 5,
};

export default function () {
    const now = Date.now();
    const vu = __VU;

    const teamName = `test_team_${vu}_${now}`;
    const users = [
        { user_id: `user1_${vu}_${now}`, username: "user1", is_active: true },
        { user_id: `user2_${vu}_${now}`, username: "user2", is_active: true },
        { user_id: `user3_${vu}_${now}`, username: "user3", is_active: true },
    ];

    const teamRes = http.post('http://localhost:8080/team/add', JSON.stringify({
        team_name: teamName,
        members: users
    }), { headers: { "Content-Type": "application/json" } });

    check(teamRes, { "TEAM ADD status 201 or 400": r => r.status === 201 || r.status === 400 });

    const prId = `pr_${vu}_${now}`;
    const prRes = http.post('http://localhost:8080/pullRequest/create', JSON.stringify({
        pull_request_id: prId,
        pull_request_name: "test_pr",
        author_id: users[0].user_id
    }), { headers: { "Content-Type": "application/json" } });

    check(prRes, { "PR CREATE status 201 or 404 or 409": r => r.status === 201 || r.status === 404 || r.status === 409 });

    if (prRes.status === 201 && prRes.json().pr.assigned_reviewers.length > 0) {
        const oldReviewer = prRes.json().pr.assigned_reviewers[0];

        const reassignRes = http.post('http://localhost:8080/pullRequest/reassign', JSON.stringify({
            pull_request_id: prId,
            old_user_id: oldReviewer
        }), { headers: { "Content-Type": "application/json" } });

        check(reassignRes, { "PR REASSIGN status 200 or 404 or 409": r => [200, 404, 409].includes(r.status) });
    }

    const mergeRes = http.post('http://localhost:8080/pullRequest/merge', JSON.stringify({
        pull_request_id: prId
    }), { headers: { "Content-Type": "application/json" } });

    check(mergeRes, { "PR MERGE status 200 or 404": r => r.status === 200 || r.status === 404 });

    const reviewRes = http.get(`http://localhost:8080/users/getReview?user_id=${users[1].user_id}`);
    check(reviewRes, { "GET REVIEW status 200": r => r.status === 200 });
}
