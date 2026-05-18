import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../services/sync_service.dart';

class OfflineIndicator extends StatefulWidget {
  const OfflineIndicator({super.key});

  @override
  State<OfflineIndicator> createState() => _OfflineIndicatorState();
}

class _OfflineIndicatorState extends State<OfflineIndicator> {
  @override
  void initState() {
    super.initState();
    // Listen to sync service
    SyncService().onConnectivityChanged = (isOnline) {
      if (mounted) setState(() {});
    };
  }

  @override
  Widget build(BuildContext context) {
    final syncService = SyncService();

    if (syncService.isOnline) {
      // Check for pending sync
      return FutureBuilder<int>(
        future: syncService.getPendingCount(),
        builder: (context, snapshot) {
          final pending = snapshot.data ?? 0;
          if (pending > 0) {
            return Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
              color: Colors.orange,
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  const Icon(Icons.sync, size: 16, color: Colors.white),
                  const SizedBox(width: 8),
                  Text(
                    '$pending pendientes',
                    style: const TextStyle(color: Colors.white, fontSize: 12),
                  ),
                ],
              ),
            );
          }
          return const SizedBox.shrink();
        },
      );
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      color: Colors.red,
      child: const Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(Icons.wifi_off, size: 16, color: Colors.white),
          SizedBox(width: 8),
          Text(
            'Sin conexión',
            style: TextStyle(color: Colors.white, fontSize: 12),
          ),
        ],
      ),
    );
  }
}

// Provider wrapper for easier integration
class ConnectivityProvider extends ChangeNotifier {
  bool _isOnline = true;
  int _pendingCount = 0;

  bool get isOnline => _isOnline;
  int get pendingCount => _pendingCount;

  Future<void> initialize() async {
    final syncService = SyncService();
    _isOnline = syncService.isOnline;
    _pendingCount = await syncService.getPendingCount();

    syncService.onConnectivityChanged = (isOnline) {
      _isOnline = isOnline;
      notifyListeners();
    };

    syncService.onSyncComplete = () async {
      _pendingCount = await syncService.getPendingCount();
      notifyListeners();
    };
  }
}