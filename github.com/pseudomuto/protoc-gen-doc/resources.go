// AUTOGENERATED CODE. DO NOT EDIT.

package gendoc

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
)

var embeddedResources = map[string]string{
	"docbook.tmpl": "H4sIAAAAAAAA/+xZ70/zNhD+3r/Cyr7t1Zu8E0OaJrdIKytoAoQo23c3ubbWHDuznUIV9X+fHIf8/tHRUtjGF0R8T+7su3uec1t88RwytAGpqOBj5wf3m4OA+yKgfDV2fn+cff3JuZiMMJGa+gwmI4SwpprB5F4KLXzB0KXw4xC4JpoKjj1rHSGUJJLwFSB3Rhmo3c68qsA3KGPOHSWJe0dC2O1K75q3IyIJci9B+ZJG5q3URcnvLShFVpnrwjmiwdhJEncWM2YdO9ZlOeKN4KuWqH1xjY0ukXtN1IwCC1S+jjVZMEBLSUIYO4SxPGAeEvuMKMVJ2IheGJB1W9uQcbGSIo6QL5gaOz+WnCOEzWIEvjE+0UCvx873jncw4pt7Pgw6q0P0GkhQXkEIS/FUXUEIA9dyO0lPiz370A553EbQj7ghC2D9kFIlW4HYq+0Re42DYL0QQe29Un9XumHw4KWG79s3ZpT/icwf4EVHm5SYjs66yD5iz8Am/f7MGyZbQ3FbGWBb/xKWJGb6D8JiE9XgJtnazyhJ6nYvBSQJ8KAjaCP3Jq0pvFqPavaxZxmR09pLCVhQuOwhJ+2vzxq4kbnjE/cOlIYAFREGOHx+Cg5/DJbnOTmU6b8QNYC4i8MFyHcWg5YuG8zROwlC099UcE0op3xV81wY/rno2LL8r1QHe5WLTtlUNAqPw9NdXY6ldGmSh+Tt7BTydqAumbP9C+TE5vvoUnIYLU9Fqg4aZQ9dY736ISNvb/Mp5CuDDbDuOY0pXwoZEtbLls9J/jnJPyf5f2qSV3i/j/hkPTIHuaH+676CePsZ3jK/b0Gvxcf4juHNBcueFQ1P+gf4Kwal0bB0PYCKBFewB/TN9Skr5QnEKctPTUmy1YaMWKpn1rmWQELKV7sdUun/fXTu3YPNfGMTdrlzF9b8ym18xNtPzVDSpdZvV+c+YUTae3vatVXuD994uu87Q5xusZ8PAY5rbxKqUaWs0m4khRbdtH65mAhtEtgNmH750mf+jWxIn/1+q9ddemHXrkRv+O96vV/f95kf4sW2xV5r7YZMNUWqGI9p81UlqqsELzMz/V2hxPDSc9/mjbiZ4gyiplG0lzdTqr2AtmZ7Qa/2O8h0viYyGg683u8kpq6dwIZw1WWrJlpVyWq5OZXUaYS9/HejvwMAAP//LC5x+GkaAAA=",
	"html.tmpl": "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",
	"markdown.tmpl": "H4sIAAAAAAAA/+RWTU/bTBC++1fMG7+HIuRwR0kOhQKqACGCekFV2SSTxJKz63rXUZF3/3u1H/au7RhSCfVSDtgzY8/H8zw7cQwPBRNsyTK4ZMtyh1QQkTIaTQhQssPpSLB8NJuckVkUxTE8kUWGwNZwwahAKnhUVQWhG4TxVZohVyqqqv/XaYY/9OtwPoXxPdmhUgk8V5W7//4pbu5PIoCqSiBdw/gOOScb5KCU8brMtVspAJvmltFNmOqqzLIwHdKVS5EA0hUkjaXLfKHlrlvD+D6swC+BlKeM9qo0AVdKg5ZkuMcMfMyU9CAqlWATO6b8HIt9uuzBWLv9lH86Ye1N4Hm+JBkp4BvJSoSn1xx119w4k712JkI7T6KjBeIl13TmhDfJgWTphk5HRbrZitFsQmBb4Ho6io06n1iun5uc5VakzftRVY0vkS+LNNeiViroxouqVdjD4FVvMnpFHMyqcb8h/CrFbKVzSjC3IA04IOGWLDADCcGbICMJif4De4W26f5Ahhjq/I5v6ScF2chW1wtJtfaJfdp0YR6vKsoWBbQnsXNc4pqUmTDcKgXOPAczdxhyejANWmXUV49IKPlIepV7ZD4Tri/35W6BxRBCfZQakIbRCs7hIcRagFlbrzWS0pRuuhHb3t+BzoYm/yX26O3qtZgks0DBbmV9gHwl6NB7HBwBuBluEOz3IAzA6IOAZmk7BA5pK4BhYHeGkBxcvP+uPl96An0ZVKjnpEVBR53+1+Y4gb4hzjsUW7aqNfqIP0vkombnEXnOKMfaHmSnS0TX7Noy/K3QDQyvXddSd/s6d7CELQnOPxcFkl1KN0oBN/cN5C6rnayf1voP5LWBtxK/d/B6HPP6Y8LR685nHEP/E0CTNc71B2XNxj0TyEHCxekpSPhK9gQkPLyKrTlg10yHYu26edBclovXIdLs9eDZ8v98PBCiadOTF6rRfPxaDEdwNoO2yzGtR2jOX56HMT1QaNvJQs91K9fFfEuKvHl620qmp6/tBujfAQAA//+3qHlyoQsAAA==",
	"scalars.json": "H4sIAAAAAAAA/9yXzW4aMRDH7zzFiFMqBZDSlEa9JZGQOOQEOUWp5GVnvW6NTewxzaqq1HfoG/ZJqt0F1gYvpChVk9zQfHg9v/nPWNx1AL53AAC6C6NJT4sFdj9BN9Uukdg9rV1KE9rSvDbMFot45MzG7XxzciY1o+H52vGFLVk8ZZEvgpyNvaBcq6jLuKRYO0aVowPw47SlxiB1X4lBYFNhYN4q8P1ZrMCwjn9dn1DeNTb13Vq0sGRGsERiT6LilAOqmU6F4n0YK8wyMROoCDJtNh5QyBmJJYJy8wSNhd8/f4HIoNDOQCZQpiAsSPEVZQGkIWdl7DppyaRDewrOItjqYiCUJWRpPwI8uHkDXKgI7iDWh+1Fe6iFIuRo4rC9FB/1leDKzUEbGInH8tcJs2DwwQmD6btDPWi0/sJ6MDw/0IPm5k0PpFY82oT4SPvhu10YWDLCC9huxsBP3+3IXvLuGPlHOLgWMbqoGt0zy/EAgaM06Y4RZQuZmERcXCPu/4jkKET2qN05EVxhCkJRPWt9mOZoEebaIGxGWhZ1Cu6OM+VMgUHuJDNQXcG+7fVoj9qPz855eN7G+VWvwEw8YhqR8aX8xgoLWflqJAWh7cNNgK4G5FbvdLZ6OYAZBJ0RKuAGGaGp484+n128vM35zEqtWEakumKJgud0CGb54B6G+WH4epft/mF/uh7f9tL7Cy29wa2UaC2f8q/Lj2vK9K08eqZfZWlHpiKFbnuCCrecfoFT4/BaMmsHIyZt/XN/twOSTa+hdsDcWQJWd36mFTGh4HY66l2sXq+01NjHXiIILifX4zEQPlJMF+GHGmKhnbdczKc2CZuft3wiZGbJDJwS5ZVj3Ooz4aQqbf98VMrfAXbDig0fpgpgJhFkmCnA4oNDNSvXafvUtNG5KggnLYTu7svjYoR2s55OaR+dqsO9i6vxtEbUue/8CQAA//8z+wC/ohEAAA==",
}

func fetchResource(name string) ([]byte, error) {
	raw, ok := embeddedResources[name]
	if !ok {
		return nil, fmt.Errorf("Could not find resource for '%s'", name)
	}

	compressed, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	buf := bytes.NewBuffer(compressed)
	
	r, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(&out, r); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
