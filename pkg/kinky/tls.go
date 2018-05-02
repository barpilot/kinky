package kinky

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "github.com/barpilot/kinky/pkg/apis/kinky/v1alpha1"
	"github.com/barpilot/kinky/pkg/tls"

	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	"github.com/operator-framework/operator-sdk/pkg/sdk/query"
)

var (
	defaultClusterDomain = "cluster.local"
	orgForTLSCert        = []string{"coreos.com"}
)

// newEtcdServerTLSSecret returns a secret containing etcd server TLS assets
func newEtcdServerTLSSecret(vr *api.Cluster, caKey *rsa.PrivateKey, caCrt *x509.Certificate) (*v1.Secret, error) {
	return newTLSSecret(vr, caKey, caCrt, "etcd server", etcdServerTLSSecretName(vr.Name),
		[]string{
			"localhost",
			fmt.Sprintf("*.%s.%s.svc", etcdNameForKinky(vr.Name), vr.Namespace),
			fmt.Sprintf("%s-client", etcdNameForKinky(vr.Name)),
			fmt.Sprintf("%s-client.%s", etcdNameForKinky(vr.Name), vr.Namespace),
			fmt.Sprintf("%s-client.%s.svc", etcdNameForKinky(vr.Name), vr.Namespace),
			fmt.Sprintf("*.%s.%s.svc.%s", etcdNameForKinky(vr.Name), vr.Namespace, defaultClusterDomain),
			fmt.Sprintf("%s-client.%s.svc.%s", etcdNameForKinky(vr.Name), vr.Namespace, defaultClusterDomain),
		},
		map[string]string{
			"key":  "server.key",
			"cert": "server.crt",
			"ca":   "server-ca.crt",
		})
}

// etcdClientTLSSecretName returns the name of etcd client TLS secret for the given vault name
func etcdClientTLSSecretName(vaultName string) string {
	return vaultName + "-etcd-client-tls"
}

// etcdServerTLSSecretName returns the name of etcd server TLS secret for the given vault name
func etcdServerTLSSecretName(vaultName string) string {
	return vaultName + "-etcd-server-tls"
}

// etcdPeerTLSSecretName returns the name of etcd peer TLS secret for the given vault name
func etcdPeerTLSSecretName(vaultName string) string {
	return vaultName + "-etcd-peer-tls"
}

// prepareEtcdTLSSecrets creates three etcd TLS secrets (client, server, peer) containing TLS assets.
// Currently we self-generate the CA, and use the self generated CA to sign all the TLS certs.
func prepareEtcdTLSSecrets(vr *api.Cluster) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("prepare TLS secrets failed: %v", err)
		}
	}()

	se := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      etcdClientTLSSecretName(vr.Name),
			Namespace: vr.Namespace,
		},
	}

	err = query.Get(se)
	if err == nil {
		return nil
	}
	if !apierrors.IsNotFound(err) {
		return err
	}

	caKey, caCrt, err := newCACert()
	if err != nil {
		return err
	}

	se, err = newEtcdClientTLSSecret(vr, caKey, caCrt)
	if err != nil {
		return err
	}
	addOwnerRefToObject(se, asOwner(vr))
	err = action.Create(se)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}

	se, err = newEtcdServerTLSSecret(vr, caKey, caCrt)
	if err != nil {
		return err
	}
	addOwnerRefToObject(se, asOwner(vr))
	err = action.Create(se)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}

	se, err = newEtcdPeerTLSSecret(vr, caKey, caCrt)
	if err != nil {
		return err
	}
	addOwnerRefToObject(se, asOwner(vr))
	err = action.Create(se)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

func newCACert() (*rsa.PrivateKey, *x509.Certificate, error) {
	key, err := tls.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}

	config := tls.CertConfig{
		CommonName:   "vault operator CA",
		Organization: orgForTLSCert,
	}

	cert, err := tls.NewSelfSignedCACertificate(config, key)
	if err != nil {
		return nil, nil, err
	}

	return key, cert, err
}

// newEtcdClientTLSSecret returns a secret containing etcd client TLS assets
func newEtcdClientTLSSecret(vr *api.Cluster, caKey *rsa.PrivateKey, caCrt *x509.Certificate) (*v1.Secret, error) {
	return newTLSSecret(vr, caKey, caCrt, "etcd client", etcdClientTLSSecretName(vr.Name), nil,
		map[string]string{
			"key":  "etcd-client.key",
			"cert": "etcd-client.crt",
			"ca":   "etcd-client-ca.crt",
		})
}

// newEtcdPeerTLSSecret returns a secret containing etcd peer TLS assets
func newEtcdPeerTLSSecret(vr *api.Cluster, caKey *rsa.PrivateKey, caCrt *x509.Certificate) (*v1.Secret, error) {
	return newTLSSecret(vr, caKey, caCrt, "etcd peer", etcdPeerTLSSecretName(vr.Name),
		[]string{
			fmt.Sprintf("*.%s.%s.svc", etcdNameForKinky(vr.Name), vr.Namespace),
			fmt.Sprintf("*.%s.%s.svc.%s", etcdNameForKinky(vr.Name), vr.Namespace, defaultClusterDomain),
		},
		map[string]string{
			"key":  "peer.key",
			"cert": "peer.crt",
			"ca":   "peer-ca.crt",
		})
}

// newTLSSecret is a common utility for creating a secret containing TLS assets.
func newTLSSecret(vr *api.Cluster, caKey *rsa.PrivateKey, caCrt *x509.Certificate, commonName, secretName string,
	addrs []string, fieldMap map[string]string) (*v1.Secret, error) {
	tc := tls.CertConfig{
		CommonName:   commonName,
		Organization: orgForTLSCert,
		AltNames:     tls.NewAltNames(addrs),
	}
	key, crt, err := newKeyAndCert(caCrt, caKey, tc)
	if err != nil {
		return nil, fmt.Errorf("new TLS secret failed: %v", err)
	}
	secret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: vr.Namespace,
			Labels:    labelsForKinky(vr.Name),
		},
		Data: map[string][]byte{
			fieldMap["key"]:  tls.EncodePrivateKeyPEM(key),
			fieldMap["cert"]: tls.EncodeCertificatePEM(crt),
			fieldMap["ca"]:   tls.EncodeCertificatePEM(caCrt),
		},
	}
	return secret, nil
}

func newKeyAndCert(caCert *x509.Certificate, caPrivKey *rsa.PrivateKey, config tls.CertConfig) (*rsa.PrivateKey, *x509.Certificate, error) {
	key, err := tls.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}
	cert, err := tls.NewSignedCertificate(config, key, caCert, caPrivKey)
	if err != nil {
		return nil, nil, err
	}
	return key, cert, nil
}