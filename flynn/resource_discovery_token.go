package flynn

import (
    "fmt"
    "path"
    "net/http"
    "net/url"
    "runtime"

    "github.com/hashicorp/terraform/helper/schema"
)

func resourceDiscoveryToken() *schema.Resource {
    return &schema.Resource{
        Create: resourceDiscoveryTokenCreate,
        Read:   resourceDiscoveryTokenRead,
        Update: resourceDiscoveryTokenUpdate,
        Delete: resourceDiscoveryTokenDelete,

        Schema: map[string]*schema.Schema{
            "token": {
                Type:     schema.TypeString,
                Computed: true,
            },
            "server": {
                Type:        schema.TypeString,
                Optional:    true,
                Default:     "https://discovery.flynn.io",
                Description: "The server used to retrieve a discovery token\nIf not set, defaults to https://discovery.flynn.io",
            },
        },
    }
}

func getDiscoveryTokenFrom(uri string) (string, error) {
    req, err := http.NewRequest("POST", uri, nil)
    if err != nil {
        return "", err
    }
    req.Header.Set("User-Agent", fmt.Sprintf("terraform-provider-flynn/%s-%s", runtime.GOOS, runtime.GOARCH))
    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }
    if res.StatusCode != http.StatusCreated {
        return "", fmt.Errorf("HTTP request error. Response code: %d", res.StatusCode)
    }

    base, err := url.Parse(uri)
    if err != nil {
        return "", err
    }
    cluster, err := url.Parse(res.Header.Get("Location"))
    if err != nil {
        return "", err
    }

    return base.ResolveReference(cluster).String(), nil
}

func resourceDiscoveryTokenCreate(d *schema.ResourceData, meta interface{}) error {
    server := d.Get("server").(string)
    uri := server + "/clusters"

    token, err := getDiscoveryTokenFrom(uri)

    if err == nil {
        d.Set("token", string(token))
        d.SetId(path.Base(token))
    } else {
        return fmt.Errorf("Error requesting discovery token: %d", err)
    }

    return resourceDiscoveryTokenRead(d, meta)
}

func resourceDiscoveryTokenRead(d *schema.ResourceData, meta interface{}) error {
    return nil
}

func resourceDiscoveryTokenUpdate(d *schema.ResourceData, meta interface{}) error {
    return resourceDiscoveryTokenRead(d, meta)
}

func resourceDiscoveryTokenDelete(d *schema.ResourceData, meta interface{}) error {
    return nil
}