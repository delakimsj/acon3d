pm.test("create workspace", function(){
    pm.collectionVariables.set("workspaceId", pm.response.json().workspace_id);
})
